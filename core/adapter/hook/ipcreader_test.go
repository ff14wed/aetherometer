package hook_test

import (
	"encoding/binary"
	"math"
	"net/url"
	"sync"

	"github.com/ff14wed/aetherometer/core/adapter/hook"
	"github.com/ff14wed/aetherometer/core/testhelpers"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("IPCReader", func() {
	var (
		ir           *hook.IPCReader
		payloadsChan chan hook.Payload

		logBuf *testhelpers.LogBuffer
		once   sync.Once

		supervisor *suture.Supervisor

		nonPDKInitZoneBlockBytes []byte
	)

	BeforeEach(func() {
		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err := zap.RegisterSink("ipcreadertest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
		logBuf.Reset()
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"ipcreadertest://"}
		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		nonPDKInitZoneBlockBytes = append(initZoneBlockBytes[:0:0], initZoneBlockBytes...)
		nonPDKInitZoneBlockBytes[2] = 0x80

		payloadsChan = make(chan hook.Payload)

		ir = hook.NewIPCReader(123, payloadsChan, logger)

		supervisor = suture.New("test-ipcreader", suture.Spec{
			Log: func(line string) {
				_, _ = GinkgoWriter.Write([]byte(line))
			},
			FailureThreshold: 1,
		})
		supervisor.ServeBackground()
		_ = supervisor.Add(ir)

	})

	AfterEach(func() {
		supervisor.Stop()
	})

	It(`logs "Running" on startup`, func() {
		Eventually(logBuf).Should(gbytes.Say("ipc-reader.*Running"))
	})

	It(`logs "Stopping..." on shutdown`, func() {
		supervisor.Stop()
		Eventually(logBuf).Should(gbytes.Say("ipc-reader.*Stopping..."))
	})

	Context("when receiving OpDebug payloads", func() {
		It("logs the message as a debug log", func() {
			payloadsChan <- hook.Payload{Op: hook.OpDebug, Channel: 123, Data: []byte("Hello")}
			Eventually(logBuf).Should(gbytes.Say("DEBUG.*ipc-reader.*.*channel.*123.*Hello"))
			Consistently(ir.SubscribeIngress()).ShouldNot(Receive())
			Consistently(ir.SubscribeEgress()).ShouldNot(Receive())
		})
	})

	Context("when receiving OpPing payloads", func() {
		It("does nothing", func() {
			payloadsChan <- hook.Payload{Op: hook.OpPing, Channel: 0, Data: []byte("Hello")}
			Consistently(ir.SubscribeIngress()).ShouldNot(Receive())
			Consistently(ir.SubscribeEgress()).ShouldNot(Receive())
		})
	})

	Context("when receiving OpExit payloads", func() {
		It("does nothing", func() {
			payloadsChan <- hook.Payload{Op: hook.OpExit, Channel: 0, Data: []byte("Hello")}
			Consistently(ir.SubscribeIngress()).ShouldNot(Receive())
			Consistently(ir.SubscribeEgress()).ShouldNot(Receive())
		})
	})

	Context("when receiving payloads of other Ops", func() {
		It("does nothing", func() {
			payloadsChan <- hook.Payload{Op: 123, Channel: 0, Data: []byte("Hello")}
			Consistently(ir.SubscribeIngress()).ShouldNot(Receive())
			Consistently(ir.SubscribeEgress()).ShouldNot(Receive())
		})
	})

	Context("when receiving OpRecv payloads", func() {
		It("decodes the ingress IPC segments and parses the blocks", func() {
			payloadsChan <- hook.Payload{Op: hook.OpRecv, Channel: 1, Data: payloadForIPCBlock(1234, 5678, datatypes.MovementOpcode, movementBlockBytes)}
			payloadsChan <- hook.Payload{Op: hook.OpRecv, Channel: 1, Data: payloadForIPCBlock(5678, 5678, datatypes.MovementOpcode, movementBlockBytes)}
			payloadsChan <- hook.Payload{Op: hook.OpRecv, Channel: 1, Data: payloadForIPCBlock(5678, 5678, datatypes.UndefinedOpcode, egressMovementBlockBytes)}
			Consistently(ir.SubscribeEgress()).ShouldNot(Receive())
			var b1, b2 *xivnet.Block
			Eventually(ir.SubscribeIngress()).Should(Receive(&b1))
			Eventually(ir.SubscribeIngress()).Should(Receive(&b2))
			Expect(b1.Data).To(Equal(expectedMovementBlockData))
			Expect(b2.Data).To(Equal(expectedMovementBlockData))

			Consistently(ir.SubscribeIngress()).ShouldNot(Receive())
		})

		It("decodes data received from different channels", func() {
			payloadsChan <- hook.Payload{Op: hook.OpRecv, Channel: 1, Data: payloadForIPCBlock(1234, 5678, datatypes.MovementOpcode, movementBlockBytes)}
			payloadsChan <- hook.Payload{Op: hook.OpRecv, Channel: 3, Data: payloadForIPCBlock(5678, 5678, datatypes.MovementOpcode, movementBlockBytes)}

			Consistently(ir.SubscribeEgress()).ShouldNot(Receive())
			var b1, b2 *xivnet.Block
			Eventually(ir.SubscribeIngress()).Should(Receive(&b1))
			Eventually(ir.SubscribeIngress()).Should(Receive(&b2))
			Expect(b1.Data).To(Equal(expectedMovementBlockData))
			Expect(b2.Data).ToNot(Equal(expectedMovementBlockData))

			Consistently(ir.SubscribeIngress()).ShouldNot(Receive())
		})

		Context("when handling data in a zone with a PDK", func() {
			BeforeEach(func() {
				payloadsChan <- hook.Payload{Op: hook.OpRecv, Channel: 1, Data: payloadForIPCBlock(2345, 5678, datatypes.InitZoneOpcode, initZoneBlockBytes)}
				payloadsChan <- hook.Payload{Op: hook.OpRecv, Channel: 1, Data: payloadForIPCBlock(2345, 5678, datatypes.CastingOpcode, castingBlockBytes)}
				Eventually(ir.SubscribeIngress()).Should(Receive())
				Eventually(ir.SubscribeIngress()).Should(Receive())
			})

			It("uses the PDK to decrypt the Lockon control block", func() {
				payloadsChan <- hook.Payload{Op: hook.OpRecv, Channel: 1, Data: payloadForIPCBlock(2345, 5678, datatypes.ControlOpcode, lockonBlockBytes)}
				var b *xivnet.Block
				Eventually(ir.SubscribeIngress()).Should(Receive(&b))
				Expect(b.Data).To(Equal(&datatypes.Control{
					Type: 0x22,
					P1:   0x2E,
				}))
			})
		})

		Context("when handling data in a zone without a PDK", func() {
			BeforeEach(func() {
				payloadsChan <- hook.Payload{Op: hook.OpRecv, Channel: 1, Data: payloadForIPCBlock(2345, 5678, datatypes.InitZoneOpcode, nonPDKInitZoneBlockBytes)}
				payloadsChan <- hook.Payload{Op: hook.OpRecv, Channel: 1, Data: payloadForIPCBlock(2345, 5678, datatypes.CastingOpcode, castingBlockBytes)}
				Eventually(ir.SubscribeIngress()).Should(Receive())
				Eventually(ir.SubscribeIngress()).Should(Receive())
			})

			It("does not do any sort of PDK handling", func() {
				payloadsChan <- hook.Payload{Op: hook.OpRecv, Channel: 1, Data: payloadForIPCBlock(2345, 5678, datatypes.ControlOpcode, lockonBlockBytes)}
				var b *xivnet.Block
				Eventually(ir.SubscribeIngress()).Should(Receive(&b))
				Expect(b.Data).To(Equal(&datatypes.Control{
					Type: 0x22,
					P1:   0x40,
				}))
			})
		})
	})

	Context("when receiving OpSend payloads", func() {
		It("decodes the egress frames and parses the blocks", func() {
			payloadsChan <- hook.Payload{Op: hook.OpSend, Channel: 1, Data: payloadForIPCBlock(1234, 5678, datatypes.UndefinedOpcode, movementBlockBytes)}
			payloadsChan <- hook.Payload{Op: hook.OpSend, Channel: 1, Data: payloadForIPCBlock(5678, 5678, datatypes.UndefinedOpcode, movementBlockBytes)}
			payloadsChan <- hook.Payload{Op: hook.OpSend, Channel: 1, Data: payloadForIPCBlock(5678, 5678, datatypes.EgressMovementOpcode, egressMovementBlockBytes)}

			Consistently(ir.SubscribeIngress()).ShouldNot(Receive())
			var b *xivnet.Block
			Eventually(ir.SubscribeEgress()).Should(Receive(&b))

			Expect(b.Data).To(Equal(expectedEgressMovementBlockData))

			Consistently(ir.SubscribeEgress()).ShouldNot(Receive())
		})
	})
})

var movementBlockBytes = []byte{
	0x12, 0x12, 0x67, 0x45, 0x01, 0x02, // HeadRotation, Direction, etc.
	0xAB, 0x89, 0xAB, 0x89, 0xAB, 0x89, // PackedPosition
	0x67, 0x45, 0x00, 0x00, // U3
}
var expectedMovementBlockData = &datatypes.Movement{
	HeadRotation:    0x12,
	Direction:       0x12,
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

var initZoneBlockBytes = []byte{
	0x12, 0x34,
	0xB2, 0x03, // TerritoryTypeID,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, // WeatherID, Bitmask
	0x00, 0x00,
	// U6 follows
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// X, Y, Z
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

var castingBlockBytes = []byte{
	0x00, 0x10, // ActionIDName
	0x00, 0x00,
	0x12, 0x10, 0x00, 0x00, // ActionID
	0x00, 0x00, 0x00, 0x00, // CastTime
	0x00, 0x00, 0x00, 0xE0, // TargetID
	0x00, 0x00, // Direction
	0x00, 0x00,
	0x00, 0x00, 0x00, 0xE0, // UnkID1
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Packed Position
	0x00, 0x00,
}

var lockonBlockBytes = []byte{
	0x22, 0x00, // Type
	0x00, 0x00,
	0x40, 0x00, 0x00, 0x00, // P1
	0x00, 0x00, 0x00, 0x00, // P2
	0x00, 0x00, 0x00, 0x00, // P3
	0x00, 0x00, 0x00, 0x00, // P4
	0x00, 0x00, 0x00, 0x00,
}

func payloadForIPCBlock(subjectID uint32, currentID uint32, opcode uint16, blockData []byte) []byte {
	header := make([]byte, 32)
	binary.LittleEndian.PutUint32(header[0:4], subjectID)
	binary.LittleEndian.PutUint32(header[4:8], currentID)
	binary.LittleEndian.PutUint64(header[8:16], 12000)
	binary.LittleEndian.PutUint16(header[16:18], 0x14)
	binary.LittleEndian.PutUint16(header[18:20], opcode)
	binary.LittleEndian.PutUint16(header[22:24], 123)
	return append(header, blockData...)
}
