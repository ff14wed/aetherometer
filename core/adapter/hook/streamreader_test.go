package hook_test

import (
	"errors"
	"io"
	"net/url"
	"sync"

	"github.com/ff14wed/aetherometer/core/adapter/hook"
	"github.com/ff14wed/aetherometer/core/adapter/hook/hookfakes"
	"github.com/ff14wed/aetherometer/core/testhelpers"
	"github.com/onsi/gomega/gbytes"
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreamReader", func() {
	type readData struct {
		data []byte
		err  error
	}

	var (
		sr *hook.StreamReader

		hookConn     *hookfakes.FakeReadCloser
		sendFakeData func(readData)

		logBuf *testhelpers.LogBuffer
		once   sync.Once

		supervisor *suture.Supervisor
	)

	BeforeEach(func() {
		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err := zap.RegisterSink("streamreadertest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
		logBuf.Reset()
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"streamreadertest://"}
		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		hookConn = new(hookfakes.FakeReadCloser)
		fakeDataChan := make(chan readData)

		sendFakeData = func(d readData) {
			fakeDataChan <- d
		}

		hookConn.ReadStub = func(p []byte) (n int, err error) {
			d, ok := <-fakeDataChan
			if !ok {
				return 0, io.EOF
			}
			copy(p, d.data)
			return len(d.data), d.err
		}
		hookConn.CloseStub = func() error {
			close(fakeDataChan)
			return nil
		}

		sr = hook.NewStreamReader(hookConn, logger)

		supervisor = suture.New("test-streamreader", suture.Spec{
			Log: func(line string) {
				_, _ = GinkgoWriter.Write([]byte(line))
			},
			FailureThreshold: 1,
		})
		supervisor.ServeBackground()
		_ = supervisor.Add(sr)
	})

	AfterEach(func() {
		supervisor.Stop()
	})

	It(`logs "Running" on startup`, func() {
		Eventually(logBuf).Should(gbytes.Say("stream-reader.*Running"))
	})

	It(`logs "Stopping..." on shutdown`, func() {
		supervisor.Stop()
		Eventually(logBuf).Should(gbytes.Say("stream-reader.*Stopping..."))
	})

	It("closes the hook connection on shutdown", func() {
		supervisor.Stop()
		Expect(hookConn.CloseCallCount()).To(Equal(1))
	})

	It("receives data on the hook connection and decodes it into envelopes", func() {
		sendFakeData(readData{
			data: append([]byte{
				20, 0, 0, 0, // Length
				200,          // Op
				210, 4, 0, 0, // Data
			}, []byte("Hello World")...),
		})

		var e hook.Envelope
		Eventually(sr.ReceivedEnvelopesListener()).Should(Receive(&e))
		Expect(e).To(Equal(hook.Envelope{
			Length: 20, Op: 200, Data: 1234, Additional: []byte("Hello World"),
		}))
	})

	Context("when the reader returns some sort of non-fatal error", func() {
		It("logs the error and continues running", func() {
			sendFakeData(readData{err: errors.New("Boom")})
			Eventually(logBuf).Should(gbytes.Say("ERROR.*stream-reader.*reading data from conn.*Boom"))

			sendFakeData(readData{
				data: []byte{9, 0, 0, 0, 200, 210, 4, 0, 0},
			})

			var e hook.Envelope
			Eventually(sr.ReceivedEnvelopesListener()).Should(Receive(&e))
			Expect(e).To(Equal(hook.Envelope{Length: 9, Op: 200, Data: 1234, Additional: []byte{}}))
		})
	})

	Context("when the reader returns an EOF error", func() {
		It("exits", func() {
			sendFakeData(readData{err: io.EOF})
			Eventually(logBuf).Should(gbytes.Say("stream-reader.*Stopping..."))
			Expect(sr.Complete()).To(BeTrue())
		})
	})
})
