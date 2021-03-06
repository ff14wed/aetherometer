package update_test

import (
	"math"
	"time"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/aetherometer/core/testassets"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/dealancer/validate.v2"
)

var _ = Describe("Casting Update", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streams   *store.Streams
		d         *datasheet.Collection
		streamID  int
		subjectID uint64
		entity    *models.Entity
		generator update.Generator

		expectedCastingInfo models.CastingInfo
	)

	BeforeEach(func() {
		*testEnv = genericSetup()
		b = testEnv.b
		streams = testEnv.streams
		d = testEnv.d
		streamID = testEnv.streamID
		subjectID = testEnv.subjectID
		entity = testEnv.entity
		generator = testEnv.generator

		d.ActionData = datasheet.ActionStore{
			Actions: map[uint32]datasheet.Action{
				203:  testassets.ExpectedActionData[203],
				4238: testassets.ExpectedActionData[4238],
			},
			Omens: testassets.ExpectedOmenData,
		}

		castingData := &datatypes.Casting{
			ActionIDName: 203,
			U1:           123,
			ActionID:     4238,
			CastTime:     1,
			TargetID:     0x5678,
			Direction:    0x8000,
			UnkID1:       123,
			U3:           123,
		}
		castingData.Position.X.SetFloat(1000)
		castingData.Position.Y.SetFloat(-1000)
		castingData.Position.Z.SetFloat(-1000)

		b.Data = castingData

		expectedCastingInfo = models.CastingInfo{
			ActionID:  203,
			StartTime: b.Time,
			CastTime:  time.Unix(1, 0),
			TargetID:  0x5678,
			Location: &models.Location{
				Orientation: 2 * math.Pi * float64(float32(0.5)),
				X:           1000,
				Y:           -1000,
				Z:           -1000,
				LastUpdated: b.Time,
			},
			ActionName:    "Skyshard",
			CastType:      2,
			EffectRange:   8,
			XAxisModifier: 0,
			Omen:          "general_1bf",
		}
	})

	It("generates an update that sets the entity's casting info based on the action name only", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())
		Expect(streamEvents).To(BeEmpty())

		Expect(entityEvents).To(ConsistOf(models.EntityEvent{
			StreamID: streamID,
			EntityID: subjectID,
			Type: models.UpdateCastingInfo{
				CastingInfo: &expectedCastingInfo,
			},
		}))

		Expect(entity.CastingInfo).To(Equal(&expectedCastingInfo))

		Expect(validate.Validate(entityEvents)).To(Succeed())
		Expect(validate.Validate(streams)).To(Succeed())
	})

	Context("when the action ID name is not found in the datasheets", func() {
		BeforeEach(func() {
			delete(d.ActionData.Actions, 203)
		})

		It("sets the action name to Unknown_X instead and returns a partial casting info", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

			expectedCastingInfo.ActionName = "Unknown_cb"
			expectedCastingInfo.CastType = 0
			expectedCastingInfo.EffectRange = 0
			expectedCastingInfo.XAxisModifier = 0
			expectedCastingInfo.Omen = ""

			Expect(entityEvents).To(ConsistOf(models.EntityEvent{
				StreamID: streamID,
				EntityID: subjectID,
				Type: models.UpdateCastingInfo{
					CastingInfo: &expectedCastingInfo,
				},
			}))

			Expect(entity.CastingInfo).To(Equal(&expectedCastingInfo))

			Expect(validate.Validate(entityEvents)).To(Succeed())
			Expect(validate.Validate(streams)).To(Succeed())
		})
	})

	Context("when the action ID is not found in the datasheets", func() {
		BeforeEach(func() {
			delete(d.ActionData.Actions, 4238)
		})

		It("should still return a full casting info", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

			Expect(entityEvents).To(ConsistOf(models.EntityEvent{
				StreamID: streamID,
				EntityID: subjectID,
				Type: models.UpdateCastingInfo{
					CastingInfo: &expectedCastingInfo,
				},
			}))

			Expect(entity.CastingInfo).To(Equal(&expectedCastingInfo))

			Expect(validate.Validate(entityEvents)).To(Succeed())
			Expect(validate.Validate(streams)).To(Succeed())
		})
	})

	entityValidationTests(testEnv, false)
})
