package stream_test

import (
	"math"
	"net/url"
	"sync"

	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/sibyl/backend/store/update"
	"github.com/ff14wed/sibyl/backend/stream"
	"github.com/ff14wed/sibyl/backend/testhelpers"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Handler", func() {
	var (
		handler    stream.Handler
		supervisor *suture.Supervisor

		streams store.Streams

		ingressChan chan *xivnet.Frame
		egressChan  chan *xivnet.Frame
		updateChan  chan store.Update
		generator   update.Generator

		logBuf *testhelpers.LogBuffer
		logger *zap.Logger

		once sync.Once
	)

	BeforeEach(func() {
		var err error
		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err = zap.RegisterSink("handlertest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
		})
		logBuf.Reset()
		Expect(err).ToNot(HaveOccurred())
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"handlertest://"}
		logger, err = zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		ingressChan = make(chan *xivnet.Frame)
		egressChan = make(chan *xivnet.Frame)
		updateChan = make(chan store.Update)
		generator = update.NewGenerator(nil)

		handler = stream.NewHandler(stream.HandlerFactoryArgs{
			StreamID:    1234,
			IngressChan: ingressChan,
			EgressChan:  egressChan,
			UpdateChan:  updateChan,
			Generator:   generator,
			Logger:      logger,
		})

		streams = store.Streams{
			Map: map[int]*models.Stream{
				0: &models.Stream{},
				1: &models.Stream{},
			},
			KeyOrder: []int{0, 1},
		}

		supervisor = suture.New("test-handler", suture.Spec{
			Log: func(line string) {
				_, _ = GinkgoWriter.Write([]byte(line))
			},
			FailureThreshold: 1,
		})
		supervisor.ServeBackground()
		_ = supervisor.Add(handler)
	})

	AfterEach(func() {
		supervisor.Stop()
	})

	It(`logs "Running" on startup`, func() {
		Eventually(logBuf).Should(gbytes.Say("stream-handler.*Running"))
	})

	It("emits an add stream update on startup", func() {
		var u store.Update
		Eventually(updateChan).Should(Receive(&u))
		streamEvents, entityEvents, err := u.ModifyStore(&streams)
		Expect(err).ToNot(HaveOccurred())
		Expect(entityEvents).To(BeEmpty())

		expectedStream := models.Stream{
			ID:          1234,
			EntitiesMap: make(map[uint64]*models.Entity),
		}

		Expect(streamEvents).To(ConsistOf(models.StreamEvent{
			StreamID: 1234,
			Type: models.AddStream{
				Stream: expectedStream,
			},
		}))

		Expect(streams.Map).To(HaveKeyWithValue(1234, &expectedStream))
		Expect(streams.KeyOrder).To(Equal([]int{0, 1, 1234}))
	})

	Context("when shutting down", func() {
		BeforeEach(func() {
			By("Properly add new stream")
			var u store.Update
			Eventually(updateChan).Should(Receive(&u))
			_, _, err := u.ModifyStore(&streams)
			Expect(err).ToNot(HaveOccurred())

			supervisor.Stop()
		})

		It(`logs "Stopping..." on shutdown`, func() {
			Eventually(logBuf).Should(gbytes.Say("stream-handler.*Stopping..."))
		})

		It("emits a remove stream update on shutdown", func() {
			var u store.Update
			Eventually(updateChan).Should(Receive(&u))
			streamEvents, entityEvents, err := u.ModifyStore(&streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(entityEvents).To(BeEmpty())
			Expect(streamEvents).To(ConsistOf(models.StreamEvent{
				StreamID: 1234,
				Type: models.RemoveStream{
					ID: 1234,
				},
			}))

			Expect(streams.Map).ToNot(HaveKey(1234))
			Expect(streams.KeyOrder).To(Equal([]int{0, 1}))
		})
	})

	Context("when an ingress frame is emitted by the stream", func() {
		BeforeEach(func() {
			Eventually(updateChan).Should(Receive())
			streams.Map[1234] = &models.Stream{
				ID:          1234,
				CharacterID: 0x12345678,
				EntitiesMap: map[uint64]*models.Entity{
					0x12345678: &models.Entity{
						Location: models.Location{},
					},
				},
			}
			Expect(streams.Map).To(HaveKey(1234))

			movementBlock := func(x, y, z float32) *xivnet.Block {
				movementData := &datatypes.Movement{
					Direction: 128,
				}
				movementData.Position.X.SetFloat(x)
				movementData.Position.Y.SetFloat(y)
				movementData.Position.Z.SetFloat(z)
				return &xivnet.Block{
					SubjectID: 0x12345678, CurrentID: 0x12345678, Data: movementData,
				}
			}

			f := &xivnet.Frame{}
			f.Blocks = append(f.Blocks, movementBlock(200, 200, -600))
			f.Blocks = append(f.Blocks, movementBlock(200, 200, -200))
			f.Blocks = append(f.Blocks, movementBlock(200, 200, 200))

			ingressChan <- f
		})

		It("converts all of the blocks to updates and applies them to the stream", func() {
			var u1, u2, u3 store.Update
			Eventually(updateChan).Should(Receive(&u1))
			Eventually(updateChan).Should(Receive(&u2))
			Eventually(updateChan).Should(Receive(&u3))

			expectLocationUpdate := func(u store.Update, x, y, z float64) {
				Expect(u).ToNot(BeNil())
				streamEvents, entityEvents, err := u.ModifyStore(&streams)
				Expect(err).ToNot(HaveOccurred())
				Expect(streamEvents).To(BeEmpty())

				expectedLocation := models.Location{
					Orientation: math.Pi, X: x, Y: y, Z: z,
				}
				Expect(entityEvents).To(ConsistOf(models.EntityEvent{
					StreamID: 1234,
					EntityID: 0x12345678,
					Type:     models.UpdateLocation{Location: expectedLocation},
				}))

				Expect(streams.Map[1234].EntitiesMap[0x12345678].Location).To(Equal(expectedLocation))
			}

			expectLocationUpdate(u1, 200, 200, -600)
			expectLocationUpdate(u2, 200, 200, -200)
			expectLocationUpdate(u3, 200, 200, 200)
		})
	})

	Context("when an egress frame is emitted by the stream", func() {
		BeforeEach(func() {
			Eventually(updateChan).Should(Receive())
			streams.Map[1234] = &models.Stream{
				ID:          1234,
				CharacterID: 0x12345678,
				EntitiesMap: map[uint64]*models.Entity{
					0x12345678: &models.Entity{
						Location: models.Location{},
					},
				},
			}
			Expect(streams.Map).To(HaveKey(1234))

			movementBlock := func(x, y, z float32) *xivnet.Block {
				movementData := &datatypes.EgressMovement{
					Direction: 0,
				}
				movementData.X = x
				movementData.Y = y
				movementData.Z = z
				return &xivnet.Block{
					SubjectID: 0x12345678, CurrentID: 0x12345678, Data: movementData,
				}
			}

			f := &xivnet.Frame{}
			f.Blocks = append(f.Blocks, movementBlock(200, 200, -600))
			f.Blocks = append(f.Blocks, movementBlock(200, 200, -200))
			f.Blocks = append(f.Blocks, movementBlock(200, 200, 200))

			egressChan <- f
		})

		It("converts all of the blocks to updates and applies them to the stream", func() {
			var u1, u2, u3 store.Update
			Eventually(updateChan).Should(Receive(&u1))
			Eventually(updateChan).Should(Receive(&u2))
			Eventually(updateChan).Should(Receive(&u3))

			expectLocationUpdate := func(u store.Update, x, y, z float64) {
				Expect(u).ToNot(BeNil())
				streamEvents, entityEvents, err := u.ModifyStore(&streams)
				Expect(err).ToNot(HaveOccurred())
				Expect(streamEvents).To(BeEmpty())

				expectedLocation := models.Location{
					Orientation: math.Pi, X: x, Y: y, Z: z,
				}
				Expect(entityEvents).To(ConsistOf(models.EntityEvent{
					StreamID: 1234,
					EntityID: 0x12345678,
					Type:     models.UpdateLocation{Location: expectedLocation},
				}))

				Expect(streams.Map[1234].EntitiesMap[0x12345678].Location).To(Equal(expectedLocation))
			}

			expectLocationUpdate(u1, 200, 200, -600)
			expectLocationUpdate(u2, 200, 200, -200)
			expectLocationUpdate(u3, 200, 200, 200)
		})
	})
})
