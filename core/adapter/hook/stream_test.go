package hook_test

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/ff14wed/aetherometer/core/adapter/hook"
	"github.com/ff14wed/aetherometer/core/adapter/hook/hookfakes"
	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/testhelpers"
	"github.com/ff14wed/xivnet/v3"
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("InitializeHook", func() {
	var (
		streamID uint32
		rpp      *hookfakes.FakeRemoteProcessProvider
		dllPath  string

		cfg hook.AdapterConfig
	)

	BeforeEach(func() {
		streamID = 1234
		rpp = new(hookfakes.FakeRemoteProcessProvider)
		dllPath = os.TempDir()

		cfg.RemoteProcessProvider = rpp
		cfg.HookConfig.DLLPath = dllPath
		cfg.HookConfig.DialRetryInterval = config.Duration(1 * time.Millisecond)
	})

	It("injects the provided DLL path into the process whose pid is the streamID", func() {
		_, err := hook.InitializeHook(streamID, cfg)
		Expect(err).ToNot(HaveOccurred())

		Expect(rpp.InjectDLLCallCount()).To(Equal(1))
		pid, path := rpp.InjectDLLArgsForCall(0)
		Expect(pid).To(Equal(streamID))
		Expect(path).To(Equal(dllPath))
	})

	It("dials the pipe with the expected pipe name", func() {
		_, err := hook.InitializeHook(streamID, cfg)
		Expect(err).ToNot(HaveOccurred())

		Expect(rpp.DialPipeCallCount()).To(Equal(1))
		pipeName, dialTimeout := rpp.DialPipeArgsForCall(0)
		Expect(pipeName).To(Equal(`\\.\pipe\deucalion-1234`))
		Expect(*dialTimeout).To(Equal(5 * time.Second))
	})

	Context("when the dll path cannot be found", func() {
		BeforeEach(func() {
			cfg.HookConfig.DLLPath = "non-existent-path"
		})

		It("returns an error immediately", func() {
			_, err := hook.InitializeHook(streamID, cfg)
			Expect(err).To(MatchError(ContainSubstring("non-existent-path")))

			Expect(rpp.InjectDLLCallCount()).To(BeZero())
			Expect(rpp.DialPipeCallCount()).To(BeZero())
		})
	})

	Context("when injecting the DLL fails", func() {
		BeforeEach(func() {
			rpp.InjectDLLReturns(errors.New("Inject Boom"))
		})

		It("returns an error immediately", func() {
			_, err := hook.InitializeHook(streamID, cfg)
			Expect(err).To(MatchError("Inject Boom"))

			Expect(rpp.DialPipeCallCount()).To(BeZero())
		})
	})

	Context("when dialing the pipe fails temporarily", func() {
		BeforeEach(func() {
			rpp.DialPipeReturns(nil, errors.New("Dial Boom"))
			rpp.DialPipeReturnsOnCall(4, nil, nil)
		})

		It("retries dialing the pipe (up to 5 times) until it succeeds", func() {
			_, err := hook.InitializeHook(streamID, cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(rpp.DialPipeCallCount()).To(Equal(5))
		})

		It("waits a retry interval before trying again", func() {
			cfg.HookConfig.DialRetryInterval = config.Duration(1 * time.Hour)
			go func() {
				_, _ = hook.InitializeHook(streamID, cfg)
			}()
			Eventually(rpp.DialPipeCallCount).Should(Equal(1))
			Consistently(rpp.DialPipeCallCount).Should(Equal(1))
		})
	})

	Context("when dialing the pipe continuously fails", func() {
		BeforeEach(func() {
			rpp.DialPipeReturns(nil, errors.New("Dial Boom"))
		})

		It("retries dialing the pipe (up to 5 times) and finally fails", func() {
			_, err := hook.InitializeHook(streamID, cfg)
			Expect(err).To(MatchError("Dial Boom"))
			Expect(rpp.DialPipeCallCount()).To(Equal(5))
		})
	})

	Describe("hook connection", func() {
		var (
			conn *hookfakes.FakeConn
		)

		BeforeEach(func() {
			conn = new(hookfakes.FakeConn)
			rpp.DialPipeReturns(conn, nil)
		})

		Describe("Write", func() {
			BeforeEach(func() {
				conn.WriteStub = func(p []byte) (int, error) {
					return len(p), nil
				}
			})

			It("writes to the connection", func() {
				hookConn, err := hook.InitializeHook(streamID, cfg)
				Expect(err).ToNot(HaveOccurred())

				n, err := hookConn.Write([]byte("Hello"))
				Expect(err).ToNot(HaveOccurred())
				Expect(n).To(Equal(5))

				Expect(conn.WriteCallCount()).To(Equal(1))
				Expect(conn.WriteArgsForCall(0)).To(Equal([]byte("Hello")))
			})

			Context("when there is an error writing", func() {
				BeforeEach(func() {
					conn.WriteStub = nil
					conn.WriteReturns(0, errors.New("Boom"))
				})

				It("returns the error", func() {
					hookConn, err := hook.InitializeHook(streamID, cfg)
					Expect(err).ToNot(HaveOccurred())

					_, err = hookConn.Write([]byte("Hello"))
					Expect(err).To(MatchError("Boom"))
				})
			})

			Context("when the pipe is closed", func() {
				BeforeEach(func() {
					conn.WriteStub = nil
					conn.WriteReturns(0, errors.New("Boom"))
					rpp.IsPipeClosedReturns(true)
				})

				It("returns an io.EOF error", func() {
					hookConn, err := hook.InitializeHook(streamID, cfg)
					Expect(err).ToNot(HaveOccurred())

					_, err = hookConn.Write([]byte("Hello"))
					Expect(err).To(Equal(io.EOF))
				})
			})
		})

		Describe("Read", func() {
			BeforeEach(func() {
				conn.ReadStub = func(p []byte) (int, error) {
					copy(p, []byte("Hello"))
					return 5, nil
				}
			})

			It("reads from the connection", func() {
				hookConn, err := hook.InitializeHook(streamID, cfg)
				Expect(err).ToNot(HaveOccurred())

				p := make([]byte, 5)
				n, err := hookConn.Read(p)
				Expect(err).ToNot(HaveOccurred())
				Expect(n).To(Equal(5))

				Expect(conn.ReadCallCount()).To(Equal(1))
				Expect(conn.ReadArgsForCall(0)).To(Equal(make([]byte, 5)))
			})

			Context("when there is an error reading", func() {
				BeforeEach(func() {
					conn.ReadStub = nil
					conn.ReadReturns(0, errors.New("Boom"))
				})

				It("returns the error", func() {
					hookConn, err := hook.InitializeHook(streamID, cfg)
					Expect(err).ToNot(HaveOccurred())

					p := make([]byte, 5)
					_, err = hookConn.Read(p)
					Expect(err).To(MatchError("Boom"))
				})
			})

			Context("when the pipe is closed", func() {
				BeforeEach(func() {
					conn.ReadStub = nil
					conn.ReadReturns(0, errors.New("Boom"))
					rpp.IsPipeClosedReturns(true)
				})

				It("returns an io.EOF error", func() {
					hookConn, err := hook.InitializeHook(streamID, cfg)
					Expect(err).ToNot(HaveOccurred())

					p := make([]byte, 5)
					_, err = hookConn.Read(p)
					Expect(err).To(Equal(io.EOF))
				})
			})
		})

		Describe("Close", func() {
			It("closes at most once", func() {
				hookConn, err := hook.InitializeHook(streamID, cfg)
				Expect(err).ToNot(HaveOccurred())

				_ = hookConn.Close()
				_ = hookConn.Close()
				_ = hookConn.Close()

				Expect(conn.CloseCallCount()).To(Equal(1))
			})

			Context("when there is an error when closing", func() {
				BeforeEach(func() {
					conn.CloseReturns(errors.New("Boom"))
				})
				It("returns the error when closing", func() {
					hookConn, err := hook.InitializeHook(streamID, cfg)
					Expect(err).ToNot(HaveOccurred())

					err = hookConn.Close()
					Expect(err).To(MatchError("Boom"))
				})

				It("subsequent calls to close don't return the error", func() {
					hookConn, err := hook.InitializeHook(streamID, cfg)
					Expect(err).ToNot(HaveOccurred())

					_ = hookConn.Close()
					err = hookConn.Close()
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})

var _ = Describe("Stream", func() {
	type readData struct {
		data []byte
		err  error
	}

	var (
		streamID uint32
		rpp      *hookfakes.FakeRemoteProcessProvider
		dllPath  string

		conn         *hookfakes.FakeConn
		fakeDataChan chan readData

		cfg hook.AdapterConfig

		hookStream hook.Stream
		logBuf     *testhelpers.LogBuffer
		once       sync.Once

		supervisor *suture.Supervisor
	)

	BeforeEach(func() {
		streamID = 1234
		rpp = new(hookfakes.FakeRemoteProcessProvider)
		dllPath = os.TempDir()

		conn = new(hookfakes.FakeConn)
		rpp.DialPipeReturns(conn, nil)
		rpp.IsPipeClosedStub = func(err error) bool {
			return err == io.EOF
		}

		fakeDataChan = make(chan readData)
		conn.WriteStub = func(d []byte) (int, error) {
			return len(d), nil
		}
		conn.ReadStub = func(p []byte) (n int, err error) {
			d, ok := <-fakeDataChan
			if !ok {
				return 0, io.EOF
			}
			copy(p, d.data)
			return len(d.data), d.err
		}
		conn.CloseStub = func() error {
			close(fakeDataChan)
			return nil
		}

		cfg.RemoteProcessProvider = rpp
		cfg.HookConfig.DLLPath = dllPath
		cfg.HookConfig.DialRetryInterval = config.Duration(1 * time.Millisecond)
		cfg.HookConfig.PingInterval = config.Duration(1 * time.Hour)

		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err := zap.RegisterSink("hookstreamtest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})

		logBuf.Reset()
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"hookstreamtest://"}

		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		hookStream = hook.NewStream(streamID, cfg, logger)
	})

	Context("when there is an error initializing the hook", func() {
		var (
			logger *zap.Logger
		)

		BeforeEach(func() {
			cfg.HookConfig.DLLPath = "non-existent-path"

			logBuf.Reset()
			zapCfg := zap.NewDevelopmentConfig()
			zapCfg.OutputPaths = []string{"hookstreamtest://"}

			var err error
			logger, err = zapCfg.Build()
			Expect(err).ToNot(HaveOccurred())
		})

		It("NewStream logs an error and returns nil (no process to start)", func() {
			hookStream = hook.NewStream(streamID, cfg, logger)
			Expect(hookStream).To(BeNil())
		})
	})

	Describe("StreamID", func() {
		It("returns the ID of the stream", func() {
			Expect(hookStream.StreamID()).To(Equal(int(streamID)))
		})
	})

	Describe("String", func() {
		It("returns the string representation of the stream", func() {
			Expect(hookStream.String()).To(Equal(fmt.Sprintf("stream-%d", int(streamID))))
		})
	})

	Context("when the stream successfully starts running", func() {
		BeforeEach(func() {
			supervisor = suture.New("test-stream", suture.Spec{
				Log: func(line string) {
					_, _ = GinkgoWriter.Write([]byte(line))
				},
				FailureThreshold: 1,
			})
			supervisor.ServeBackground()
			_ = supervisor.Add(hookStream)
		})

		AfterEach(func() {
			supervisor.Stop()
		})

		It(`logs "Running" for each subprocess on startup`, func() {
			Eventually(logBuf).Should(gbytes.Say("stream-1234"))
			// These subprocesses can start up in any order
			Eventually(logBuf.Buffer().Contents).Should(ContainSubstring("stream-sender"))
			Eventually(logBuf.Buffer().Contents).Should(ContainSubstring("stream-pinger"))
			Eventually(logBuf.Buffer().Contents).Should(ContainSubstring("stream-reader"))
			Eventually(logBuf.Buffer().Contents).Should(ContainSubstring("ipc-reader"))

			Eventually(logBuf).Should(gbytes.Say("Running"))
			Eventually(logBuf).Should(gbytes.Say("Running"))
			Eventually(logBuf).Should(gbytes.Say("Running"))
			Eventually(logBuf).Should(gbytes.Say("Running"))
		})

		It("initializes the hook on startup", func() {
			Expect(rpp.InjectDLLCallCount()).To(Equal(1))
			Expect(rpp.DialPipeCallCount()).To(Equal(1))
			pipeName, dialTimeout := rpp.DialPipeArgsForCall(0)
			Expect(pipeName).To(Equal(`\\.\pipe\deucalion-1234`))
			Expect(*dialTimeout).To(Equal(5 * time.Second))
		})

		It(`logs "Stopping..." for each subprocess on shutdown`, func() {
			Eventually(logBuf).Should(gbytes.Say("Running"))
			Eventually(logBuf).Should(gbytes.Say("Running"))
			Eventually(logBuf).Should(gbytes.Say("Running"))
			Eventually(logBuf).Should(gbytes.Say("Running"))

			supervisor.Stop()
			Eventually(logBuf).Should(gbytes.Say("stream-1234"))
			// These subprocesses can stop in any order
			Eventually(logBuf.Buffer().Contents).Should(ContainSubstring("stream-sender"))
			Eventually(logBuf.Buffer().Contents).Should(ContainSubstring("stream-pinger"))
			Eventually(logBuf.Buffer().Contents).Should(ContainSubstring("stream-reader"))
			Eventually(logBuf.Buffer().Contents).Should(ContainSubstring("ipc-reader"))

			Eventually(logBuf).Should(gbytes.Say("Stopping..."))
			Eventually(logBuf).Should(gbytes.Say("Stopping..."))
			Eventually(logBuf).Should(gbytes.Say("Stopping..."))
			Eventually(logBuf).Should(gbytes.Say("Stopping..."))
		})

		It("closes the connection with the hook on shutdown", func() {
			Eventually(logBuf).Should(gbytes.Say("Running"))
			Eventually(logBuf).Should(gbytes.Say("Running"))
			Eventually(logBuf).Should(gbytes.Say("Running"))
			Eventually(logBuf).Should(gbytes.Say("Running"))

			supervisor.Stop()
			Eventually(conn.CloseCallCount).Should(Equal(1))
		})

		Describe("SubscribeIngress", func() {
			It("returns ingress blocks when data is received from the hook", func() {
				Consistently(hookStream.SubscribeIngress).ShouldNot(Receive())

				fakeDataChan <- readData{data: hook.Payload{
					Op:   hook.OpRecv,
					Data: someBlockPacket,
				}.Encode()}

				Consistently(hookStream.SubscribeEgress).ShouldNot(Receive())
				var b *xivnet.Block
				Eventually(hookStream.SubscribeIngress).Should(Receive(&b))
				Expect(b.SubjectID).To(Equal(uint32(123456789)))
				Expect(b.CurrentID).To(Equal(uint32(123456789)))
				Expect(b.Time).To(Equal(time.Unix(1516257780, 175000000)))
				Expect(b.Opcode).To(Equal(uint16(0xFFFF)))
			})
		})

		Describe("SubscribeEgress", func() {
			It("returns egress blocks when data is received from the hook", func() {
				Consistently(hookStream.SubscribeEgress).ShouldNot(Receive())

				fakeDataChan <- readData{data: hook.Payload{
					Op:   hook.OpSend,
					Data: someBlockPacket,
				}.Encode()}

				Consistently(hookStream.SubscribeIngress).ShouldNot(Receive())
				var b *xivnet.Block
				Eventually(hookStream.SubscribeEgress).Should(Receive(&b))
				Expect(b.SubjectID).To(Equal(uint32(123456789)))
				Expect(b.CurrentID).To(Equal(uint32(123456789)))
				Expect(b.Time).To(Equal(time.Unix(1516257780, 175000000)))
				Expect(b.Opcode).To(Equal(uint16(0xFFFF)))
			})
		})

		Describe("SendRequest", func() {
			It("sends a JSON-encoded request as an payload on the connection", func() {
				// Byte arrays can be represented as base64 in JSON
				resp, err := hookStream.SendRequest(
					[]byte(`{"Op": 123, "Channel": 456, "Data": "BwgJAA=="}`),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp).To(Equal([]byte("OK")))

				Eventually(conn.WriteCallCount).Should(Equal(1))
				Expect(conn.WriteArgsForCall(0)).To(Equal([]byte{
					13, 0, 0, 0, 123, 200, 1, 0, 0, 7, 8, 9, 0,
				}))
			})

			Context("when the request cannot be unmarshaled", func() {
				It("returns an error", func() {
					resp, err := hookStream.SendRequest([]byte(`"bar"`))
					Expect(err).To(MatchError("cannot unmarshal data to payload: json: cannot unmarshal string into Go value of type hook.Payload"))
					Expect(resp).To(BeNil())
				})
			})
		})
	})
})

var someBlockPacket = []byte{
	0x15, 0xCD, 0x5B, 0x07, // SubjectID
	0x15, 0xCD, 0x5B, 0x07, // CurrentID
	0xcf, 0xa1, 0x01, 0x08, 0x61, 0x01, 0x00, 0x00, // Time
	0x14, 0x00, 0xFF, 0xFF, // Reserved and Opcode
	0x00, 0x00, 0x22, 0x00, // Padding and Route
	0x3f, 0xe0, 0x89, 0x58, // Time
	0x00, 0x00, 0x00, 0x00, // Pad3

	// BlockData
	0x00, 0x00, 0xcd, 0xe1, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00,
}

type alreadyInjectedErr struct{}

func (alreadyInjectedErr) Error() string              { return "Boom" }
func (alreadyInjectedErr) IsDLLAlreadyInjectedError() {}
