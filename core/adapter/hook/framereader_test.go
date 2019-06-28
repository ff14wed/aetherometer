package hook_test

import (
	"bufio"
	"errors"
	"io"
	"math"
	"net/url"
	"sync"
	"time"

	"github.com/ff14wed/aetherometer/core/adapter/hook"
	"github.com/ff14wed/aetherometer/core/message"
	"github.com/ff14wed/aetherometer/core/testhelpers"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("FrameReader", func() {
	var (
		fr            *hook.FrameReader
		envelopesChan chan hook.Envelope

		testFrames     map[string]*xivnet.Frame
		parsedFrames   map[string]*xivnet.Frame
		unparsedFrames map[string]*xivnet.Frame

		logBuf *testhelpers.LogBuffer
		once   sync.Once

		supervisor *suture.Supervisor
	)

	BeforeEach(func() {
		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err := zap.RegisterSink("framereadertest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
		logBuf.Reset()
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"framereadertest://"}
		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		testFrames = map[string]*xivnet.Frame{
			"1": &xivnet.Frame{
				Time: time.Unix(12, 0),
				Blocks: []*xivnet.Block{
					&xivnet.Block{
						Length: 123, SubjectID: 1234, CurrentID: 5678,
						IPCHeader: xivnet.IPCHeader{Opcode: datatypes.MovementOpcode},
						Data:      xivnet.GenericBlockDataFromBytes(movementBlockBytes),
					},
					&xivnet.Block{
						Length: 456, SubjectID: 5678, CurrentID: 5678,
						IPCHeader: xivnet.IPCHeader{Opcode: datatypes.MovementOpcode},
						Data:      xivnet.GenericBlockDataFromBytes(movementBlockBytes),
					},
				},
			},
			"2": &xivnet.Frame{
				Time: time.Unix(12, 0),
				Blocks: []*xivnet.Block{
					&xivnet.Block{
						Length: 789, SubjectID: 2345, CurrentID: 5678,
						IPCHeader: xivnet.IPCHeader{Opcode: datatypes.EgressMovementOpcode},
						Data:      xivnet.GenericBlockDataFromBytes(egressMovementBlockBytes),
					},
				},
			},
		}

		parsedFrames = map[string]*xivnet.Frame{
			"1": transformBlocks(testFrames["1"], func(b xivnet.Block) xivnet.Block {
				b.Time = time.Unix(12, 0)
				b.Data = expectedMovementBlockData
				return b
			}),
			"2": transformBlocks(testFrames["2"], func(b xivnet.Block) xivnet.Block {
				b.Time = time.Unix(12, 0)
				b.Data = expectedEgressMovementBlockData
				return b
			}),
		}

		unparsedFrames = map[string]*xivnet.Frame{
			"1": transformBlocks(testFrames["1"], func(b xivnet.Block) xivnet.Block {
				b.Time = time.Unix(12, 0)
				return b
			}),
			"2": transformBlocks(testFrames["2"], func(b xivnet.Block) xivnet.Block {
				b.Time = time.Unix(12, 0)
				return b
			}),
		}

		envelopesChan = make(chan hook.Envelope)

		fr = hook.NewFrameReader(123, envelopesChan, newTestFrameDecoder(testFrames), logger)

		supervisor = suture.New("test-framereader", suture.Spec{
			Log: func(line string) {
				_, _ = GinkgoWriter.Write([]byte(line))
			},
			FailureThreshold: 1,
		})
		supervisor.ServeBackground()
		_ = supervisor.Add(fr)

	})

	AfterEach(func() {
		supervisor.Stop()
	})

	It(`logs "Running" on startup`, func() {
		Eventually(logBuf).Should(gbytes.Say("frame-reader.*Running"))
	})

	It(`logs "Stopping..." on shutdown`, func() {
		supervisor.Stop()
		Eventually(logBuf).Should(gbytes.Say("frame-reader.*Stopping..."))
	})

	Context("when receiving OpDebug envelopes", func() {
		It("logs the message as a debug log", func() {
			envelopesChan <- hook.Envelope{Op: hook.OpDebug, Data: 123, Additional: []byte("Hello")}
			Eventually(logBuf).Should(gbytes.Say("DEBUG.*frame-reader.*.*data.*123.*Hello"))
		})
	})

	Context("when receiving OpPing envelopes", func() {
		It("does nothing", func() {
			envelopesChan <- hook.Envelope{Op: hook.OpPing, Data: 0, Additional: []byte("Hello")}
			Consistently(fr.SubscribeIngress()).ShouldNot(Receive())
			Consistently(fr.SubscribeEgress()).ShouldNot(Receive())
		})
	})

	Context("when receiving OpExit envelopes", func() {
		It("does nothing", func() {
			envelopesChan <- hook.Envelope{Op: hook.OpExit, Data: 0, Additional: []byte("Hello")}
			Consistently(fr.SubscribeIngress()).ShouldNot(Receive())
			Consistently(fr.SubscribeEgress()).ShouldNot(Receive())
		})
	})

	Context("when receiving OpRecv envelopes", func() {
		It("decodes the ingress frames and parses the blocks", func() {
			envelopesChan <- hook.Envelope{Op: hook.OpRecv, Data: 1, Additional: []byte("12")}
			Consistently(fr.SubscribeEgress()).ShouldNot(Receive())
			var f1, f2 *xivnet.Frame
			Eventually(fr.SubscribeIngress()).Should(Receive(&f1))
			Eventually(fr.SubscribeIngress()).Should(Receive(&f2))
			Expect(f1).To(Equal(parsedFrames["1"]))
			Expect(f2).To(Equal(unparsedFrames["2"]))

			Consistently(fr.SubscribeIngress()).ShouldNot(Receive())
		})

		It("decodes data received from different contexts", func() {
			envelopesChan <- hook.Envelope{Op: hook.OpRecv, Data: 1, Additional: []byte("1")}
			envelopesChan <- hook.Envelope{Op: hook.OpRecv, Data: 2, Additional: []byte("2")}

			Consistently(fr.SubscribeEgress()).ShouldNot(Receive())
			var f1, f2 *xivnet.Frame
			Eventually(fr.SubscribeIngress()).Should(Receive(&f1))
			Eventually(fr.SubscribeIngress()).Should(Receive(&f2))
			Expect(f1).To(Equal(parsedFrames["1"]))
			Expect(f2).To(Equal(unparsedFrames["2"]))

			Consistently(fr.SubscribeIngress()).ShouldNot(Receive())
		})

		Context("when there is an error reading the next frame", func() {
			BeforeEach(func() {
				envelopesChan <- hook.Envelope{Op: hook.OpRecv, Data: 1, Additional: []byte("3")}
			})
			It("logs a non-fatal error", func() {
				Eventually(logBuf).Should(gbytes.Say(`ERROR.*frame-reader.*Error reading next frame.*\{"context": 1, "data": "33", "isEgress": false, "error": "invalid data"\}`))
			})

			It("continues decoding frames", func() {
				envelopesChan <- hook.Envelope{Op: hook.OpRecv, Data: 1, Additional: []byte("1")}

				Consistently(fr.SubscribeEgress()).ShouldNot(Receive())
				var f *xivnet.Frame
				Eventually(fr.SubscribeIngress()).Should(Receive(&f))
				Expect(f).To(Equal(parsedFrames["1"]))

				Consistently(fr.SubscribeIngress()).ShouldNot(Receive())
			})
		})

		Context("when there is an error parsing some block in a decoded frame", func() {
			BeforeEach(func() {
				testFrames["1"].Blocks[0].Data = xivnet.GenericBlockDataFromBytes(movementBlockBytes[1:])
				unparsedFrames["1"].Blocks[0].Data = xivnet.GenericBlockDataFromBytes(movementBlockBytes[1:])
				envelopesChan <- hook.Envelope{Op: hook.OpRecv, Data: 1, Additional: []byte("1")}
			})

			It("logs a non-fatal error", func() {
				Eventually(logBuf).Should(gbytes.Say(`ERROR.*frame-reader.*Error unmarshaling block.*length mismatch`))
			})

			It("parses the rest of the blocks", func() {
				var f *xivnet.Frame
				Eventually(fr.SubscribeIngress()).Should(Receive(&f))
				Expect(f).To(Equal(&xivnet.Frame{
					Time: time.Unix(12, 0),
					Blocks: []*xivnet.Block{
						unparsedFrames["1"].Blocks[0],
						parsedFrames["1"].Blocks[1],
					},
				}))

				Consistently(fr.SubscribeIngress()).ShouldNot(Receive())
			})
		})
	})

	Context("when receiving OpSend envelopes", func() {
		It("decodes the egress frames and parses the blocks", func() {
			envelopesChan <- hook.Envelope{Op: hook.OpSend, Data: 1, Additional: []byte("12")}
			Consistently(fr.SubscribeIngress()).ShouldNot(Receive())
			var f1, f2 *xivnet.Frame
			Eventually(fr.SubscribeEgress()).Should(Receive(&f1))
			Eventually(fr.SubscribeEgress()).Should(Receive(&f2))
			Expect(f1).To(Equal(unparsedFrames["1"]))
			Expect(f2).To(Equal(parsedFrames["2"]))

			Consistently(fr.SubscribeEgress()).ShouldNot(Receive())
		})
	})
})

var movementBlockBytes = []byte{
	0x12, 0x12, 0x67, 0x45, 0x01, 0x02, // Direction, U1, U2
	0xAB, 0x89, 0xAB, 0x89, 0xAB, 0x89, // PackedPosition
	0x67, 0x45, 0x00, 0x00, // U3
}
var expectedMovementBlockData = &datatypes.Movement{
	Direction:       0x12,
	HeadRotation:    0x12,
	AnimationType:   0x67,
	AnimationState:  0x45,
	AnimationSpeed:  0x01,
	UnknownRotation: 0x02,
	Position:        datatypes.PackedPosition{X: 0x89AB, Y: 0x89AB, Z: 0x89AB},
	U3:              0x4567,
}

var egressMovementBlockBytes = []byte{
	219, 15, 73, 64, // Direction
	0x67, 0x45, 0x00, 0x00, // U1
	0, 0, 250, 67, // X
	0, 0, 22, 68, // Y
	0, 0, 47, 68, // Z
	0xAB, 0x89, 0x00, 0x00, // U2
}

var expectedEgressMovementBlockData = &datatypes.EgressMovement{
	Direction: math.Pi,
	U1:        0x4567,
	X:         500,
	Y:         600,
	Z:         700,
	U2:        0x89AB,
}

var newTestFrameDecoder = func(frames map[string]*xivnet.Frame) func(io.Reader) message.FrameDecoder {
	return func(r io.Reader) message.FrameDecoder {
		reader := bufio.NewReader(r)
		return &testFrameDecoder{buf: reader, testFrames: frames}
	}
}

type testFrameDecoder struct {
	buf        *bufio.Reader
	testFrames map[string]*xivnet.Frame
}

func (d *testFrameDecoder) NextFrame() (*xivnet.Frame, error) {
	token, err := d.buf.Peek(1)
	if err != nil {
		return nil, xivnet.EOFError{}
	}
	key := token[0]
	_, _ = d.buf.Discard(1)
	if f, ok := d.testFrames[string(key)]; ok {
		return f, nil
	}
	return nil, errors.New("invalid data")
}

func (d *testFrameDecoder) DiscardDataUntilValid() {}

func transformBlocks(f *xivnet.Frame, t func(xivnet.Block) xivnet.Block) *xivnet.Frame {
	fCopy := *f
	var newBlocks []*xivnet.Block
	for _, b := range fCopy.Blocks {
		newB := t(*b)
		newBlocks = append(newBlocks, &newB)
	}
	fCopy.Blocks = newBlocks
	return &fCopy
}
