package hook_test

import (
	"net/url"
	"sync"

	"github.com/ff14wed/aetherometer/core/adapter/hook"
	"github.com/ff14wed/aetherometer/core/adapter/hook/hookfakes"
	"github.com/ff14wed/aetherometer/core/stream"
	"github.com/ff14wed/aetherometer/core/testhelpers"
	"github.com/onsi/gomega/gbytes"
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager", func() {
	var (
		mgr *hook.Manager

		cfg              hook.AdapterConfig
		streamUp         chan stream.Provider
		streamDown       chan int
		addProcEventChan chan uint32
		remProcEventChan chan uint32

		fakeStream *hookfakes.FakeStream

		logBuf *testhelpers.LogBuffer
		once   sync.Once

		supervisor *suture.Supervisor
	)

	BeforeEach(func() {
		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err := zap.RegisterSink("hookmanagertest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
		logBuf.Reset()
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"hookmanagertest://"}
		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		streamUp = make(chan stream.Provider, 10)
		streamDown = make(chan int, 10)
		addProcEventChan = make(chan uint32, 10)
		remProcEventChan = make(chan uint32, 10)

		cfg = hook.AdapterConfig{
			StreamUp:   streamUp,
			StreamDown: streamDown,
		}

		fakeStream = new(hookfakes.FakeStream)

		closeFakeStream := make(chan struct{})
		fakeStream.ServeStub = func() {
			<-closeFakeStream
		}

		fakeStream.StringReturns("fake-stream")

		var closeOnce sync.Once
		fakeStream.StopStub = func() {
			closeOnce.Do(func() {
				close(closeFakeStream)
			})
		}

		streamBuilder := func(uint32) hook.Stream {
			return fakeStream
		}

		supervisor = suture.New("test-hookmanager", suture.Spec{
			Log: func(line string) {
				_, _ = GinkgoWriter.Write([]byte(line))
			},
			FailureThreshold: 1,
		})

		mgr = hook.NewManager(
			cfg,
			addProcEventChan,
			remProcEventChan,
			streamBuilder,
			supervisor,
			logger,
		)

		supervisor.ServeBackground()
		_ = supervisor.Add(mgr)
	})

	AfterEach(func() {
		supervisor.Stop()
	})

	It(`logs "Running" on startup`, func() {
		Eventually(logBuf).Should(gbytes.Say("hook-manager.*Running"))
	})

	It(`logs "Stopping..." on shutdown`, func() {
		supervisor.Stop()
		Eventually(logBuf).Should(gbytes.Say("hook-manager.*Stopping..."))
	})

	Context("when a process is added", func() {
		It("creates the stream and sends it on the StreamUp channel", func() {
			addProcEventChan <- 1234
			Eventually(streamUp).Should(Receive(Equal(fakeStream)))
		})

		It("creates the stream and runs it in the provided supervisor", func() {
			addProcEventChan <- 1234

			Eventually(fakeStream.ServeCallCount).Should(Equal(1))
			Consistently(fakeStream.StopCallCount).Should(BeZero())
		})
	})

	Context("when a process is removed", func() {
		BeforeEach(func() {
			addProcEventChan <- 1234

			Consistently(streamDown).ShouldNot(Receive())
			Consistently(fakeStream.StopCallCount).Should(BeZero())

			Eventually(streamUp).Should(Receive())
		})

		It("sends the ID of the stream shutting down to the StreamDown channel", func() {
			remProcEventChan <- 1234
			Eventually(streamDown).Should(Receive(Equal(1234)))
		})

		It("closes the Stream associated with the stream ID", func() {
			remProcEventChan <- 1234

			Eventually(fakeStream.StopCallCount).Should(Equal(1))
		})

		Context("when a process that doesn't exist is removed", func() {
			It("logs the error and continues operation", func() {
				remProcEventChan <- 4567
				Consistently(streamDown).ShouldNot(Receive())
				Consistently(fakeStream.StopCallCount).Should(BeZero())

				Eventually(logBuf).Should(gbytes.Say("hook-manager.*Error removing process group"))

				addProcEventChan <- 4567
				Eventually(streamUp).Should(Receive(Equal(fakeStream)))
				Eventually(fakeStream.ServeCallCount).Should(BeNumerically(">", 1))
			})
		})
	})

})
