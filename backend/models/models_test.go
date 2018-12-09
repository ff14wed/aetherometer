package models_test

import (
	"github.com/ff14wed/sibyl/backend/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Models", func() {
	Describe("DB", func() {
		var db *models.DB
		BeforeEach(func() {
			db = &models.DB{
				StreamsMap: map[int]models.Stream{
					1234: models.Stream{
						Pid: 1234,
						EntitiesMap: map[int]models.Entity{
							1: models.Entity{ID: 1, Name: "FooBar"},
							2: models.Entity{ID: 2, Name: "Baah"},
						},
						EntitiesKeys: []int{1, 2},
					},
					5678: models.Stream{Pid: 5678},
				},
				StreamKeys: []int{1234, 5678},
			}
		})
		Describe("Streams", func() {
			It("returns all the streams found in the database", func() {
				Expect(db.Streams()).To(Equal([]models.Stream{
					db.StreamsMap[1234],
					db.StreamsMap[5678],
				}))
			})
		})
		Describe("Stream", func() {
			It("returns the requested stream from the database", func() {
				Expect(db.Stream(5678)).To(Equal(models.Stream{Pid: 5678}))
			})

			It("returns an error if the requested stream does not exist", func() {
				_, err := db.Stream(2345)
				Expect(err).To(MatchError("stream ID 2345 not found"))
			})
		})
		Describe("Entity", func() {
			It("returns the requested entity from the database", func() {
				Expect(db.Entity(1234, 1)).To(Equal(models.Entity{ID: 1, Name: "FooBar"}))
			})

			It("returns an error if the requested stream does not exist", func() {
				_, err := db.Entity(2345, 1)
				Expect(err).To(MatchError("stream ID 2345 not found"))
			})

			It("returns an error if the requested entity does not exist", func() {
				_, err := db.Entity(1234, 3)
				Expect(err).To(MatchError("stream id 1234: entity ID 3 not found"))
			})
		})
	})
	Describe("Stream", func() {
		var stream *models.Stream
		BeforeEach(func() {
			stream = &models.Stream{
				Pid: 1234,
				EntitiesMap: map[int]models.Entity{
					1: models.Entity{ID: 1, Name: "FooBar"},
					2: models.Entity{ID: 2, Name: "Baah"},
				},
				EntitiesKeys: []int{1, 2},
			}
		})
		Describe("Entity", func() {
			It("returns the requested entity from the stream", func() {
				Expect(stream.Entity(1)).To(Equal(models.Entity{ID: 1, Name: "FooBar"}))
			})

			It("returns an error if the requested entity does not exist", func() {
				_, err := stream.Entity(3)
				Expect(err).To(MatchError("entity ID 3 not found"))
			})
		})
		Describe("Entities", func() {
			It("returns all entities found on the stream", func() {
				Expect(stream.Entities()).To(Equal([]models.Entity{
					models.Entity{ID: 1, Name: "FooBar"},
					models.Entity{ID: 2, Name: "Baah"},
				}))
			})
		})
	})
})
