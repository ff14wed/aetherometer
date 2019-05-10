package hook_test

import (
	"errors"
	"net/url"
	"sync"

	"github.com/ff14wed/sibyl/backend/adapter/hook"
	"github.com/ff14wed/sibyl/backend/adapter/hook/hookfakes"
	"github.com/ff14wed/sibyl/backend/testhelpers"
	"github.com/onsi/gomega/gbytes"
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreamSender", func() {
	var (
		ss *hook.StreamSender

		hookConn *hookfakes.FakeWriteCloser

		logBuf *testhelpers.LogBuffer
		once   sync.Once

		supervisor *suture.Supervisor
	)

	BeforeEach(func() {
		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err := zap.RegisterSink("streamsendertest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
		logBuf.Reset()
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"streamsendertest://"}
		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		hookConn = new(hookfakes.FakeWriteCloser)
		ss = hook.NewStreamSender(hookConn, logger)

		supervisor = suture.New("test-streamsender", suture.Spec{
			Log: func(line string) {
				_, _ = GinkgoWriter.Write([]byte(line))
			},
			FailureThreshold: 1,
		})
		supervisor.ServeBackground()
		_ = supervisor.Add(ss)

		hookConn.WriteStub = func(d []byte) (int, error) {
			return len(d), nil
		}
	})

	AfterEach(func() {
		supervisor.Stop()
	})

	It(`logs "Running" on startup`, func() {
		Eventually(logBuf).Should(gbytes.Say("stream-sender.*Running"))
	})

	It(`logs "Stopping..." on shutdown`, func() {
		supervisor.Stop()
		Eventually(logBuf).Should(gbytes.Say("stream-sender.*Stopping..."))
	})

	It("encodes and sends envelopes along the hook connection", func() {
		ss.Send(200, 1234, []byte("Hello World"))
		Eventually(hookConn.WriteCallCount).Should(Equal(1))
		Expect(hookConn.WriteArgsForCall(0)).To(Equal(append([]byte{
			20, 0, 0, 0, // Length
			200,          // Op
			210, 4, 0, 0, // Data
		}, []byte("Hello World")...)))
	})

	Context("when the writer wrote fewer bytes than expected", func() {
		BeforeEach(func() {
			hookConn.WriteStub = nil
			hookConn.WriteReturns(0, nil)
		})

		It("logs an error, but continues running", func() {
			ss.Send(200, 1234, []byte("Hello World"))
			Eventually(hookConn.WriteCallCount).Should(Equal(1))
			Eventually(logBuf).Should(gbytes.Say("ERROR.*stream-sender.*writing bytes to conn.*wrote less than the envelope length"))

			ss.Send(200, 1234, []byte("Hello World"))
			Eventually(hookConn.WriteCallCount).Should(Equal(2))
		})
	})

	Context("when the writer returns an error writing", func() {
		BeforeEach(func() {
			hookConn.WriteStub = nil
			hookConn.WriteReturns(0, errors.New("Boom"))
		})

		It("logs an error, but continues running", func() {
			ss.Send(200, 1234, []byte("Hello World"))
			Eventually(hookConn.WriteCallCount).Should(Equal(1))
			Eventually(logBuf).Should(gbytes.Say("ERROR.*stream-sender.*writing bytes to conn.*Boom"))

			ss.Send(200, 1234, []byte("Hello World"))
			Eventually(hookConn.WriteCallCount).Should(Equal(2))
		})
	})
})
