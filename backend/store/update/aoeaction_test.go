package update_test

import (
	"time"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/sibyl/backend/store/update"
	"github.com/ff14wed/xivnet/v2"
	"github.com/ff14wed/xivnet/v2/datatypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
)

var _ = Describe("AoEAction8 Update", func() {
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
			456: datasheet.Action{ID: 456, Name: "Foo"},
		}

		b.Data = &datatypes.AoEAction8{
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
				NumAffected:       2,
			},
			EffectsList: [8]datatypes.ActionEffects{
				datatypes.ActionEffects{
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
				datatypes.ActionEffects{
					datatypes.ActionEffect{
						Type:        3,
						HitSeverity: 1,
						P3:          22,
						Percentage:  50,
						Multiplier:  0,
						Flags:       0x40,
						Damage:      123,
					},
				},
				// And then some garbage data
				datatypes.ActionEffects{
					datatypes.ActionEffect{
						Type:        0xFF,
						HitSeverity: 0xFF,
						P3:          0xFF,
						Percentage:  0xFF,
						Multiplier:  0xFF,
						Flags:       0xFF,
						Damage:      0xFF,
					},
				},
			},
			Targets: [8]uint64{0xABCDEF01, 0xABCDEF02, 0x12345678},
		}

		matchExpectedAction = gstruct.MatchAllFields(gstruct.Fields{
			"TargetID":          Equal(subjectID),
			"Name":              Equal("Foo"),
			"GlobalCounter":     Equal(1),
			"AnimationLockTime": BeNumerically("~", 0.5),
			"HiddenAnimation":   Equal(2),
			"Location": gstruct.MatchAllFields(gstruct.Fields{
				"X":           Equal(float64(0)),
				"Y":           Equal(float64(0)),
				"Z":           Equal(float64(0)),
				"Orientation": BeNumerically("~", 5.445427316156579),
				"LastUpdated": Equal(time.Unix(12, 0)),
			}),
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
				models.ActionEffect{
					TargetID:        0xABCDEF02,
					Type:            3,
					HitSeverity:     1,
					Param:           22,
					BonusPercent:    50,
					ValueMultiplier: 0,
					Flags:           0x40,
					Value:           123,
				},
			),
			"EffectFlags": Equal(0),
			"UseTime":     Equal(time.Unix(12, 0)),
		})
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
		Expect(*entity.LastAction).To(matchExpectedAction)
	})
})

var _ = Describe("AoEAction16 Update", func() {
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
			456: datasheet.Action{ID: 456, Name: "Foo"},
		}

		b.Data = &datatypes.AoEAction16{
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
				NumAffected:       2,
			},
			EffectsList: [16]datatypes.ActionEffects{
				datatypes.ActionEffects{
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
				datatypes.ActionEffects{
					datatypes.ActionEffect{
						Type:        3,
						HitSeverity: 1,
						P3:          22,
						Percentage:  50,
						Multiplier:  0,
						Flags:       0x40,
						Damage:      123,
					},
				},
				// And then some garbage data
				datatypes.ActionEffects{
					datatypes.ActionEffect{
						Type:        0xFF,
						HitSeverity: 0xFF,
						P3:          0xFF,
						Percentage:  0xFF,
						Multiplier:  0xFF,
						Flags:       0xFF,
						Damage:      0xFF,
					},
				},
			},
			Targets: [16]uint64{0xABCDEF01, 0xABCDEF02, 0x12345678},
		}

		matchExpectedAction = gstruct.MatchAllFields(gstruct.Fields{
			"TargetID":          Equal(subjectID),
			"Name":              Equal("Foo"),
			"GlobalCounter":     Equal(1),
			"AnimationLockTime": BeNumerically("~", 0.5),
			"HiddenAnimation":   Equal(2),
			"Location": gstruct.MatchAllFields(gstruct.Fields{
				"X":           Equal(float64(0)),
				"Y":           Equal(float64(0)),
				"Z":           Equal(float64(0)),
				"Orientation": BeNumerically("~", 5.445427316156579),
				"LastUpdated": Equal(time.Unix(12, 0)),
			}),
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
				models.ActionEffect{
					TargetID:        0xABCDEF02,
					Type:            3,
					HitSeverity:     1,
					Param:           22,
					BonusPercent:    50,
					ValueMultiplier: 0,
					Flags:           0x40,
					Value:           123,
				},
			),
			"EffectFlags": Equal(0),
			"UseTime":     Equal(time.Unix(12, 0)),
		})
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
		Expect(*entity.LastAction).To(matchExpectedAction)
	})
})

var _ = Describe("AoEAction24 Update", func() {
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
			456: datasheet.Action{ID: 456, Name: "Foo"},
		}

		b.Data = &datatypes.AoEAction24{
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
				NumAffected:       2,
			},
			EffectsList: [24]datatypes.ActionEffects{
				datatypes.ActionEffects{
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
				datatypes.ActionEffects{
					datatypes.ActionEffect{
						Type:        3,
						HitSeverity: 1,
						P3:          22,
						Percentage:  50,
						Multiplier:  0,
						Flags:       0x40,
						Damage:      123,
					},
				},
				// And then some garbage data
				datatypes.ActionEffects{
					datatypes.ActionEffect{
						Type:        0xFF,
						HitSeverity: 0xFF,
						P3:          0xFF,
						Percentage:  0xFF,
						Multiplier:  0xFF,
						Flags:       0xFF,
						Damage:      0xFF,
					},
				},
			},
			Targets: [24]uint64{0xABCDEF01, 0xABCDEF02, 0x12345678},
		}

		matchExpectedAction = gstruct.MatchAllFields(gstruct.Fields{
			"TargetID":          Equal(subjectID),
			"Name":              Equal("Foo"),
			"GlobalCounter":     Equal(1),
			"AnimationLockTime": BeNumerically("~", 0.5),
			"HiddenAnimation":   Equal(2),
			"Location": gstruct.MatchAllFields(gstruct.Fields{
				"X":           Equal(float64(0)),
				"Y":           Equal(float64(0)),
				"Z":           Equal(float64(0)),
				"Orientation": BeNumerically("~", 5.445427316156579),
				"LastUpdated": Equal(time.Unix(12, 0)),
			}),
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
				models.ActionEffect{
					TargetID:        0xABCDEF02,
					Type:            3,
					HitSeverity:     1,
					Param:           22,
					BonusPercent:    50,
					ValueMultiplier: 0,
					Flags:           0x40,
					Value:           123,
				},
			),
			"EffectFlags": Equal(0),
			"UseTime":     Equal(time.Unix(12, 0)),
		})
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
		Expect(*entity.LastAction).To(matchExpectedAction)
	})
})

var _ = Describe("AoEAction32 Update", func() {
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
			456: datasheet.Action{ID: 456, Name: "Foo"},
		}

		b.Data = &datatypes.AoEAction32{
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
				NumAffected:       2,
			},
			EffectsList: [32]datatypes.ActionEffects{
				datatypes.ActionEffects{
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
				datatypes.ActionEffects{
					datatypes.ActionEffect{
						Type:        3,
						HitSeverity: 1,
						P3:          22,
						Percentage:  50,
						Multiplier:  0,
						Flags:       0x40,
						Damage:      123,
					},
				},
				// And then some garbage data
				datatypes.ActionEffects{
					datatypes.ActionEffect{
						Type:        0xFF,
						HitSeverity: 0xFF,
						P3:          0xFF,
						Percentage:  0xFF,
						Multiplier:  0xFF,
						Flags:       0xFF,
						Damage:      0xFF,
					},
				},
			},
			Targets: [32]uint64{0xABCDEF01, 0xABCDEF02, 0x12345678},
		}

		matchExpectedAction = gstruct.MatchAllFields(gstruct.Fields{
			"TargetID":          Equal(subjectID),
			"Name":              Equal("Foo"),
			"GlobalCounter":     Equal(1),
			"AnimationLockTime": BeNumerically("~", 0.5),
			"HiddenAnimation":   Equal(2),
			"Location": gstruct.MatchAllFields(gstruct.Fields{
				"X":           Equal(float64(0)),
				"Y":           Equal(float64(0)),
				"Z":           Equal(float64(0)),
				"Orientation": BeNumerically("~", 5.445427316156579),
				"LastUpdated": Equal(time.Unix(12, 0)),
			}),
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
				models.ActionEffect{
					TargetID:        0xABCDEF02,
					Type:            3,
					HitSeverity:     1,
					Param:           22,
					BonusPercent:    50,
					ValueMultiplier: 0,
					Flags:           0x40,
					Value:           123,
				},
			),
			"EffectFlags": Equal(0),
			"UseTime":     Equal(time.Unix(12, 0)),
		})
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
		Expect(*entity.LastAction).To(matchExpectedAction)
	})
})
