package hook_test

import (
	"net/url"
	"sync"
	"time"

	"github.com/ff14wed/aetherometer/core/adapter/hook"
	"github.com/ff14wed/aetherometer/core/adapter/hook/hookfakes"
	"github.com/ff14wed/aetherometer/core/testhelpers"
	"github.com/onsi/gomega/gbytes"
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreamPinger", func() {
	var (
		sp  *hook.StreamPinger
		hds *hookfakes.FakeHookDataSender

		logBuf *testhelpers.LogBuffer
		once   sync.Once

		supervisor *suture.Supervisor
	)

	BeforeEach(func() {
		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err := zap.RegisterSink("streampingertest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
		logBuf.Reset()
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"streampingertest://"}
		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		hds = new(hookfakes.FakeHookDataSender)
		sp = hook.NewStreamPinger(hds, 100*time.Millisecond, logger)

		supervisor = suture.New("test-streampinger", suture.Spec{
			Log: func(line string) {
				_, _ = GinkgoWriter.Write([]byte(line))
			},
			FailureThreshold: 1,
		})
		supervisor.ServeBackground()
		_ = supervisor.Add(sp)
	})

	AfterEach(func() {
		supervisor.Stop()
	})

	It(`logs "Running" on startup`, func() {
		Eventually(logBuf).Should(gbytes.Say("stream-pinger.*Running"))
	})

	It(`logs "Stopping..." on shutdown`, func() {
		supervisor.Stop()
		Eventually(logBuf).Should(gbytes.Say("stream-pinger.*Stopping..."))
	})

	It("periodically sends a ping to the HookDataSender", func() {
		for i := 0; i < 10; i++ {
			Eventually(hds.SendCallCount).Should(BeNumerically(">", i+1))
			op, _, _ := hds.SendArgsForCall(i)
			Expect(op).To(Equal(byte(hook.OpPing)))
		}
	})

	It("doesn't send any more pings after shutdown", func() {
		supervisor.Stop()
		Eventually(logBuf).Should(gbytes.Say("stream-pinger.*Stopping..."))

		sendCallCount := hds.SendCallCount()
		Consistently(hds.SendCallCount).Should(Equal(sendCallCount))
	})
})
