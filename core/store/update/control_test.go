package update_test

import (
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/dealancer/validate.v2"
)

var _ = Describe("Control Update", func() {
	Describe("type 0xF", func() {
		var (
			testEnv = new(testVars)

			b         *xivnet.Block
			streams   *store.Streams
			streamID  int
			subjectID uint64
			entity    *models.Entity
			generator update.Generator
		)

		BeforeEach(func() {
			*testEnv = genericSetup()
			b = testEnv.b
			streams = testEnv.streams
			streamID = testEnv.streamID
			subjectID = testEnv.subjectID
			entity = testEnv.entity
			generator = testEnv.generator

			notifyData := &datatypes.Control{
				Type: 0xF,
			}

			b.Data = notifyData
		})

		It("does nothing", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).To(BeNil())
		})

		Context("when P1 is 538", func() {
			BeforeEach(func() {
				notifyData := &datatypes.Control{
					Type: 0xF,
					P1:   538,
				}

				b.Data = notifyData
			})

			It("generates an update that nulls out the entity's casting info", func() {
				u := generator.Generate(streamID, false, b)
				Expect(u).ToNot(BeNil())
				streamEvents, entityEvents, err := u.ModifyStore(streams)
				Expect(err).ToNot(HaveOccurred())
				Expect(streamEvents).To(BeEmpty())

				Expect(entityEvents).To(ConsistOf(models.EntityEvent{
					StreamID: streamID,
					EntityID: subjectID,
					Type: models.UpdateCastingInfo{
						CastingInfo: nil,
					},
				}))

				Expect(entity.CastingInfo).To(BeNil())

				Expect(validate.Validate(entityEvents)).To(Succeed())
				Expect(validate.Validate(streams)).To(Succeed())
			})

			entityValidationTests(testEnv, false)
		})
	})

	Describe("type 0x22", func() {
		var (
			testEnv = new(testVars)

			b         *xivnet.Block
			streams   *store.Streams
			streamID  int
			subjectID uint64
			entity    *models.Entity
			generator update.Generator

			expectedLockonMarker int
		)

		BeforeEach(func() {
			*testEnv = genericSetup()
			b = testEnv.b
			streams = testEnv.streams
			streamID = testEnv.streamID
			subjectID = testEnv.subjectID
			entity = testEnv.entity
			generator = testEnv.generator

			expectedLockonMarker = 123

			notifyData := &datatypes.Control{
				Type: 0x22,
				P1:   uint32(expectedLockonMarker),
			}

			b.Data = notifyData
		})

		It("generates an update that sets the entity's lockon marker", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

			Expect(entityEvents).To(ConsistOf(models.EntityEvent{
				StreamID: streamID,
				EntityID: subjectID,
				Type: models.UpdateLockonMarker{
					LockonMarker: expectedLockonMarker,
				},
			}))

			Expect(entity.LockonMarker).To(Equal(expectedLockonMarker))

			Expect(validate.Validate(entityEvents)).To(Succeed())
			Expect(validate.Validate(streams)).To(Succeed())
		})

		entityValidationTests(testEnv, false)
	})
})
