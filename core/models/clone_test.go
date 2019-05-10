package models_test

import (
	"github.com/ff14wed/aetherometer/core/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stream Clone", func() {
	var stream *models.Stream

	BeforeEach(func() {
		bNPCInfoName := "Bar"
		bNPCInfoSize := 1.5

		stream = &models.Stream{
			ID:          1234,
			CharacterID: 4567,
			Place: models.Place{
				MapID: 20,
				Maps: []models.MapInfo{
					models.MapInfo{ID: "some-map"},
				},
			},
			Enmity: models.Enmity{
				TargetHateRanking: []models.HateRanking{
					models.HateRanking{ActorID: 2345, Hate: 1000},
				},
				NearbyEnemyHate: []models.HateEntry{
					models.HateEntry{EnemyID: 3456, HatePercent: 99},
				},
			},

			CraftingInfo: &models.CraftingInfo{
				StepNum: 900,
			},
			EntitiesMap: map[uint64]*models.Entity{
				1: &models.Entity{
					ID: 1, Index: 2, Name: "FooBar",
					BNPCInfo: &models.NPCInfo{
						Name: &bNPCInfoName,
						Size: &bNPCInfoSize,
					},
					LastAction: &models.Action{
						TargetID: 5678,
						Effects: []models.ActionEffect{
							models.ActionEffect{TargetID: 2345},
						},
					},
					CastingInfo: &models.CastingInfo{ActionID: 100},
					Statuses: []*models.Status{
						&models.Status{ID: 50},
					},
				},
			},
		}
	})

	It("produces an identical copy of the original struct", func() {
		streamClone := stream.Clone()
		Expect(streamClone).To(Equal(*stream))
	})

	DescribeTable("changes on the clone of the stream should not affect the original copy",
		func(modifier func(*models.Stream)) {
			streamClone := stream.Clone()
			modifier(&streamClone)
			Expect(streamClone).ToNot(Equal(*stream))
		},
		Entry("Stream", func(s *models.Stream) {
			s.CharacterID = 5678
		}),
		Entry("stream.Place.Maps", func(s *models.Stream) {
			s.Place.Maps[0].ID = "foo-map"
		}),
		Entry("stream.Enmity.TargetHateRanking", func(s *models.Stream) {
			s.Enmity.TargetHateRanking[0].Hate = 1001
		}),
		Entry("stream.Enmity.NearbyEnemyHate", func(s *models.Stream) {
			s.Enmity.NearbyEnemyHate[0].HatePercent = 100
		}),
		Entry("stream.CraftingInfo", func(s *models.Stream) {
			s.CraftingInfo.StepNum = 200
		}),
		Entry("stream.EntitiesMap", func(s *models.Stream) {
			s.EntitiesMap[2] = &models.Entity{ID: 2, Name: "Baah", Index: 1}
		}),
		Entry("Entity", func(s *models.Stream) {
			s.EntitiesMap[1].Name = "BarFoo"
		}),
		Entry("entity.BNPCInfo", func(s *models.Stream) {
			s.EntitiesMap[1].BNPCInfo.NameID = 1
		}),
		Entry("entity.BNPCInfo.Name", func(s *models.Stream) {
			newName := "Baz"
			s.EntitiesMap[1].BNPCInfo.Name = &newName
		}),
		Entry("entity.BNPCInfo.Size", func(s *models.Stream) {
			newSize := 2.0
			s.EntitiesMap[1].BNPCInfo.Size = &newSize
		}),
		Entry("entity.LastAction", func(s *models.Stream) {
			s.EntitiesMap[1].LastAction.TargetID = 6789
		}),
		Entry("entity.LastAction.Effects", func(s *models.Stream) {
			s.EntitiesMap[1].LastAction.Effects[0].TargetID = 3456
		}),
		Entry("entity.CastingInfo", func(s *models.Stream) {
			s.EntitiesMap[1].CastingInfo.ActionID = 101
		}),
		Entry("entity.Statuses", func(s *models.Stream) {
			s.EntitiesMap[1].Statuses[0].ID = 51
		}),
	)

	Describe("with a minimal stream", func() {
		BeforeEach(func() {
			stream = &models.Stream{
				ID: 1234,
				EntitiesMap: map[uint64]*models.Entity{
					1: &models.Entity{ID: 1, Index: 2, Name: "FooBar"},
					2: &models.Entity{ID: 1, Index: 2, Name: "FooBar", BNPCInfo: &models.NPCInfo{}},
					3: nil,
				},
			}
		})

		It("produces an identical copy of the original struct", func() {
			streamClone := stream.Clone()
			Expect(streamClone).To(Equal(*stream))
		})
	})

	Describe("with an empty stream", func() {
		BeforeEach(func() {
			stream = &models.Stream{ID: 1234}
		})

		It("produces an identical copy of the original struct", func() {
			streamClone := stream.Clone()
			Expect(streamClone).To(Equal(*stream))
		})
	})
})
