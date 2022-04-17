package update_test

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	"gopkg.in/dealancer/validate.v2"
)

var _ = Describe("Action Update", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streams   *store.Streams
		d         *datasheet.Collection
		streamID  int
		subjectID uint64
		entity    *models.Entity
		generator update.Generator

		matchExpectedAction types.GomegaMatcher
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
				456: {Key: 456, Name: "Foo"},
			},
		}

		b.Data = &datatypes.Action{
			ActionHeader: datatypes.ActionHeader{
				TargetID:          uint32(subjectID),
				ActionIDName:      456,
				GlobalCounter:     1,
				AnimationLockTime: 0.5,
				HiddenAnimation:   2,
				Direction:         0xDDDD,
				ActionID:          123,
				Variation:         3,
				EffectDisplayType: 4,
				NumAffected:       1,
			},
			Effects: datatypes.ActionEffects{
				datatypes.ActionEffect{
					Type:        3,
					HitSeverity: 1,
					P3:          22,
					Percentage:  50,
					Multiplier:  0,
					Flags:       0x40,
					Damage:      123,
				},
				datatypes.ActionEffect{
					Type:        4,
					HitSeverity: 0,
					P3:          23,
					Percentage:  75,
					Multiplier:  0,
					Flags:       0x41,
					Damage:      456,
				},
			},
			TargetID2:   0xABCDEF01,
			EffectFlags: 5,
		}

		matchExpectedAction = gstruct.PointTo(gstruct.MatchAllFields(gstruct.Fields{
			"TargetID":          Equal(subjectID),
			"Name":              Equal("Foo"),
			"GlobalCounter":     Equal(1),
			"AnimationLockTime": BeNumerically("~", 0.5),
			"HiddenAnimation":   Equal(2),
			"Location": gstruct.PointTo(gstruct.MatchAllFields(gstruct.Fields{
				"X":           Equal(float64(0)),
				"Y":           Equal(float64(0)),
				"Z":           Equal(float64(0)),
				"Orientation": BeNumerically("~", 5.445427316156579),
				"LastUpdated": Equal(b.Time),
			})),
			"ID":                Equal(123),
			"Variation":         Equal(3),
			"EffectDisplayType": Equal(4),
			"IsAoE":             BeFalse(),
			"Effects": ConsistOf(
				models.ActionEffect{
					TargetID:        0xABCDEF01,
					Type:            3,
					HitSeverity:     1,
					Param:           22,
					BonusPercent:    50,
					ValueMultiplier: 0,
					Flags:           0x40,
					Value:           123,
				},
				models.ActionEffect{
					TargetID:        0xABCDEF01,
					Type:            4,
					HitSeverity:     0,
					Param:           23,
					BonusPercent:    75,
					ValueMultiplier: 0,
					Flags:           0x41,
					Value:           456,
				},
			),
			"EffectFlags": Equal(5),
			"UseTime":     Equal(b.Time),
		}))
	})

	It("generates an update that sets the entity's last action", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())
		Expect(streamEvents).To(BeEmpty())

		Expect(entityEvents).To(HaveLen(1))
		Expect(entityEvents[0].StreamID).To(Equal(streamID))
		Expect(entityEvents[0].EntityID).To(Equal(subjectID))
		eventType, assignable := entityEvents[0].Type.(models.UpdateLastAction)
		Expect(assignable).To(BeTrue())
		Expect(eventType.Action).To(matchExpectedAction)

		Expect(entity.LastAction).ToNot(BeNil())
		Expect(entity.LastAction).To(matchExpectedAction)

		Expect(validate.Validate(entityEvents)).To(Succeed())
		Expect(validate.Validate(streams)).To(Succeed())
	})

	Context("when the action ID name is not found in the datasheets", func() {
		BeforeEach(func() {
			delete(d.ActionData.Actions, 456)
		})

		It("sets the action name to Unknown_X instead", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

			Expect(entityEvents).To(HaveLen(1))
			Expect(entityEvents[0].StreamID).To(Equal(streamID))
			Expect(entityEvents[0].EntityID).To(Equal(subjectID))
			eventType, assignable := entityEvents[0].Type.(models.UpdateLastAction)
			Expect(assignable).To(BeTrue())
			Expect(eventType.Action.Name).To(Equal("Unknown_1c8"))

			Expect(entity.LastAction).ToNot(BeNil())
			Expect(entity.LastAction.Name).To(Equal("Unknown_1c8"))

			Expect(validate.Validate(entityEvents)).To(Succeed())
			Expect(validate.Validate(streams)).To(Succeed())
		})
	})

	Context("when a casting info is present on the entity", func() {
		BeforeEach(func() {
			streams.Map[streamID].EntitiesMap[subjectID].CastingInfo =
				&models.CastingInfo{ActionID: 1234, ActionName: "Bar"}
		})

		It("generates update that sets the entity's last action and removes the casting info", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

			Expect(entityEvents).To(HaveLen(2))
			Expect(entityEvents).To(ContainElement(
				models.EntityEvent{
					StreamID: streamID,
					EntityID: subjectID,
					Type: models.UpdateCastingInfo{
						CastingInfo: nil,
					},
				},
			))
			Expect(entity.LastAction).ToNot(BeNil())
			Expect(entity.LastAction).To(matchExpectedAction)
			Expect(entity.CastingInfo).To(BeNil())

			Expect(validate.Validate(entityEvents)).To(Succeed())
			Expect(validate.Validate(streams)).To(Succeed())
		})
	})

	entityValidationTests(testEnv, false)
})
