package update_test

import (
	"time"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store/update"
	"github.com/ff14wed/xivnet"
	"github.com/ff14wed/xivnet/datatypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Action Update", func() {
	var (
		b  *xivnet.Block
		db *models.DB
		d  *datasheet.Collection

		stream    int
		subjectID uint64
		entity    *models.Entity

		matchExpectedAction types.GomegaMatcher
	)

	BeforeEach(func() {
		stream = 1234
		subjectID = 0x12345678

		entity = &models.Entity{}

		db = &models.DB{
			StreamsMap: map[int]*models.Stream{
				stream: &models.Stream{
					Pid: stream,
					EntitiesMap: map[uint64]*models.Entity{
						subjectID: entity,
					},
				},
			},
		}

		d = new(datasheet.Collection)
		d.ActionData = datasheet.ActionStore{
			456: datasheet.Action{ID: 456, Name: "Foo"},
		}

		b = &xivnet.Block{
			Length: 1234,
			Header: xivnet.BlockHeader{
				SubjectID: uint32(subjectID),
				CurrentID: 0x9ABCDEF0,
				Opcode:    1234,
				Time:      time.Unix(12, 0),
			},
			Data: &datatypes.Action{
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
			},
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
			),
			"EffectFlags": Equal(5),
			"UseTime":     Equal(time.Unix(12, 0)),
		})
	})

	It("generates an update that sets the entity's last action", func() {
		g := update.NewGenerator(stream, false, d)
		u := g.Generate(b)
		Expect(u).ToNot(BeNil())
		streamEvents, entityEvents, err := u.ModifyDB(db)
		Expect(err).ToNot(HaveOccurred())
		Expect(streamEvents).To(BeEmpty())

		Expect(entityEvents).To(HaveLen(1))
		Expect(entityEvents[0].StreamID).To(Equal(stream))
		Expect(entityEvents[0].EntityID).To(Equal(subjectID))
		eventType, assignable := entityEvents[0].Type.(models.UpdateLastAction)
		Expect(assignable).To(BeTrue())
		Expect(eventType.Action).To(matchExpectedAction)

		Expect(entity.LastAction).ToNot(BeNil())
		Expect(*entity.LastAction).To(matchExpectedAction)
	})

	It("errors when the stream doesn't exist", func() {
		g := update.NewGenerator(1000, false, d)
		u := g.Generate(b)
		Expect(u).ToNot(BeNil())

		streamEvents, entityEvents, err := u.ModifyDB(db)
		Expect(err).To(MatchError(update.ErrorStreamNotFound))
		Expect(streamEvents).To(BeEmpty())
		Expect(entityEvents).To(BeEmpty())
	})

	It("errors when the entity doesn't exist", func() {
		b.Header.SubjectID = 0x9ABCDEF0

		g := update.NewGenerator(1234, false, d)
		u := g.Generate(b)
		Expect(u).ToNot(BeNil())

		streamEvents, entityEvents, err := u.ModifyDB(db)
		Expect(err).To(MatchError(update.ErrorEntityNotFound))
		Expect(streamEvents).To(BeEmpty())
		Expect(entityEvents).To(BeEmpty())
	})

	Context("when the action ID name is not found in the datasheets", func() {
		BeforeEach(func() {
			delete(d.ActionData, 456)
		})

		It("sets the action name to Unknown_X instead", func() {
			g := update.NewGenerator(stream, false, d)
			u := g.Generate(b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyDB(db)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

			Expect(entityEvents).To(HaveLen(1))
			Expect(entityEvents[0].StreamID).To(Equal(stream))
			Expect(entityEvents[0].EntityID).To(Equal(subjectID))
			eventType, assignable := entityEvents[0].Type.(models.UpdateLastAction)
			Expect(assignable).To(BeTrue())
			Expect(eventType.Action.Name).To(Equal("Unknown_1c8"))

			Expect(entity.LastAction).ToNot(BeNil())
			Expect(entity.LastAction.Name).To(Equal("Unknown_1c8"))

		})
	})

	Context("when a casting info is present on the entity", func() {
		BeforeEach(func() {
			db.StreamsMap[stream].EntitiesMap[subjectID].CastingInfo =
				&models.CastingInfo{ActionID: 1234, ActionName: "Bar"}
		})

		It("generates update that sets the entity's last action and removes the casting info", func() {
			g := update.NewGenerator(stream, false, d)
			u := g.Generate(b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyDB(db)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

			Expect(entityEvents).To(HaveLen(2))
			Expect(entityEvents).To(ContainElement(
				models.EntityEvent{
					StreamID: stream,
					EntityID: subjectID,
					Type: models.UpdateCastingInfo{
						CastingInfo: nil,
					},
				},
			))
			Expect(entity.LastAction).ToNot(BeNil())
			Expect(*entity.LastAction).To(matchExpectedAction)
			Expect(entity.CastingInfo).To(BeNil())
		})
	})
})
