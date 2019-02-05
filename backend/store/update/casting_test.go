package update_test

import (
	"math"
	"time"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/sibyl/backend/store/update"
	"github.com/ff14wed/sibyl/backend/testassets"
	"github.com/ff14wed/xivnet"
	"github.com/ff14wed/xivnet/datatypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			203:  testassets.ExpectedActionData[203],
			4238: testassets.ExpectedActionData[4238],
		}

		castingData := &datatypes.Casting{
			ActionIDName: 203,
			U1:           123,
			ActionID:     4238,
			CastTime:     1,
			TargetID:     0x5678,
			Direction:    math.Pi,
			UnkID1:       123,
			U3:           123,
		}
		castingData.Position.X.SetFloat(1000)
		castingData.Position.Y.SetFloat(-1000)
		castingData.Position.Z.SetFloat(-1000)

		b.Data = castingData

		expectedCastingInfo = models.CastingInfo{
			ActionID:  4238,
			StartTime: time.Unix(12, 0),
			CastTime:  time.Unix(1, 0),
			TargetID:  0x5678,
			Location: models.Location{
				Orientation: float64(float32(math.Pi)),
				X:           1000,
				Y:           -1000,
				Z:           -1000,
			},
			ActionName:    "Skyshard",
			CastType:      4,
			EffectRange:   30,
			XAxisModifier: 4,
			Omen:          "general02f",
		}
	})

	It("generates an update that sets the entity's casting info", func() {
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
	})

	Context("when the action ID name is not found in the datasheets", func() {
		BeforeEach(func() {
			delete(d.ActionData, 203)
		})

		It("sets the action name to Unknown_X instead", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

			expectedCastingInfo.ActionName = "Unknown_cb"

			Expect(entityEvents).To(ConsistOf(models.EntityEvent{
				StreamID: streamID,
				EntityID: subjectID,
				Type: models.UpdateCastingInfo{
					CastingInfo: &expectedCastingInfo,
				},
			}))

			Expect(entity.CastingInfo).To(Equal(&expectedCastingInfo))
		})
	})

	Context("when the action ID is not found in the datasheets", func() {
		BeforeEach(func() {
			delete(d.ActionData, 4238)
		})

		It("sets a partially blank casting info", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

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
		})
	})

	entityValidationTests(testEnv)
})
