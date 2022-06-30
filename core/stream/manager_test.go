package stream_test

import (
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/aetherometer/core/stream"
	"github.com/ff14wed/aetherometer/core/stream/streamfakes"
	"github.com/ff14wed/aetherometer/core/testhelpers"
	"github.com/ff14wed/xivnet/v3"
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

type FakeHandler struct {
	serveCalled uint32
	stopCalled  uint32
	stop        chan struct{}
}

func (f *FakeHandler) Serve() {
	atomic.StoreUint32(&f.serveCalled, 1)
	<-f.stop
}

func (f *FakeHandler) ServeCalled() bool {
	return atomic.LoadUint32(&f.serveCalled) == 1
}

func (f *FakeHandler) Stop() {
	atomic.StoreUint32(&f.stopCalled, 1)
	close(f.stop)
}

func (f *FakeHandler) StopCalled() bool {
	return atomic.LoadUint32(&f.stopCalled) == 1
}

var _ = Describe("Manager", func() {
	var (
		manager    *stream.Manager
		supervisor *suture.Supervisor

		generator        update.Generator
		updateChan       chan<- store.Update
		streamSupervisor *suture.Supervisor

		handlerFactoryArgs stream.HandlerFactoryArgs
		fakeHandler        *FakeHandler

		logBuf *testhelpers.LogBuffer
		once   sync.Once
	)

	BeforeEach(func() {
		var err error
		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err := zap.RegisterSink("managertest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
		logBuf.Reset()
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"managertest://"}
		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		generator = update.NewGenerator(nil)
		updateChan = make(chan store.Update)
		fakeHandler = &FakeHandler{stop: make(chan struct{})}
		handlerFactory := func(args stream.HandlerFactoryArgs) stream.Handler {
			handlerFactoryArgs = args
			return fakeHandler
		}
		streamSupervisor = suture.New("stream-supervisor", suture.Spec{
			Log: func(line string) {
				logger.Named("stream-supervisor").Info(line)
			},
		})

		manager = stream.NewManager(
			generator,
			updateChan,
			streamSupervisor,
			handlerFactory,
			logger,
		)

		supervisor = suture.New("test-manager", suture.Spec{
			Log: func(line string) {
				_, _ = GinkgoWriter.Write([]byte(line))
			},
			FailureThreshold: 1,
		})
		supervisor.ServeBackground()
		_ = supervisor.Add(streamSupervisor)
		_ = supervisor.Add(manager)
	})

	AfterEach(func() {
		supervisor.Stop()
	})

	It(`logs "Running" on startup`, func() {
		Eventually(logBuf).Should(gbytes.Say("stream-manager.*Running"))
	})

	It(`logs "Stopping..." on shutdown`, func() {
		supervisor.Stop()
		Eventually(logBuf).Should(gbytes.Say("stream-manager.*Stopping..."))
	})

	It("logs an error when attempting to remove a non-existent stream", func() {
		manager.StreamDown() <- 1234
		Eventually(logBuf).Should(gbytes.Say("Error removing stream.*1234"))
	})

	Context("when a new stream is created", func() {
		var (
			fakeProvider *streamfakes.FakeProvider
			ingressChan  <-chan *xivnet.Frame
			egressChan   <-chan *xivnet.Frame
		)

		BeforeEach(func() {
			fakeProvider = new(streamfakes.FakeProvider)
			fakeProvider.StreamIDReturns(1234)
			ingressChan = make(chan *xivnet.Frame)
			egressChan = make(chan *xivnet.Frame)
			fakeProvider.SubscribeIngressReturns(ingressChan)
			fakeProvider.SubscribeEgressReturns(egressChan)
			fakeProvider.SendRequestReturns([]byte("ack"), nil)
			manager.StreamUp() <- fakeProvider

			Eventually(fakeHandler.ServeCalled).Should(BeTrue())
			Eventually(fakeProvider.SubscribeEgressCallCount()).Should(BeNumerically(">", 0))
		})

		It("creates a new Handler for the stream", func() {
			Expect(handlerFactoryArgs.StreamID).To(Equal(1234))
			Expect(handlerFactoryArgs.IngressChan).To(Equal(ingressChan))
			Expect(handlerFactoryArgs.EgressChan).To(Equal(egressChan))
			Expect(handlerFactoryArgs.UpdateChan).To(Equal(updateChan))
			Expect(handlerFactoryArgs.Generator).To(Equal(generator))
		})

		Describe("SendRequest", func() {
			It("handles requests for the stream", func() {
				resp, err := manager.SendRequest(1234, []byte("foo"))
				Expect(err).ToNot(HaveOccurred())
				Expect(resp).To(Equal([]byte("ack")))
			})

			It("errors when the stream doesn't exist", func() {
				_, err := manager.SendRequest(5678, nil)
				Expect(err).To(MatchError("stream provider 5678 not found"))
			})
		})

		Context("when the stream is closed", func() {
			BeforeEach(func() {
				Consistently(fakeHandler.StopCalled).Should(BeFalse())
				manager.StreamDown() <- 1234
			})

			It("stops the Handler for the stream", func() {
				Eventually(fakeHandler.StopCalled).Should(BeTrue())
			})

			It("can no longer handle requests for the stream", func() {
				Eventually(fakeHandler.StopCalled).Should(BeTrue())
				_, err := manager.SendRequest(1234, nil)
				Expect(err).To(MatchError("stream provider 1234 not found"))
			})
		})
	})
})
