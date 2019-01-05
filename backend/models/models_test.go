package models_test

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/models/modelsfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Models", func() {
	Describe("DB", func() {
		var (
			db                    *models.DB
			fakeStreamEventSource *modelsfakes.FakeStreamEventSource
			fakeEntityEventSource *modelsfakes.FakeEntityEventSource
		)

		BeforeEach(func() {
			fakeStreamEventSource = new(modelsfakes.FakeStreamEventSource)
			fakeEntityEventSource = new(modelsfakes.FakeEntityEventSource)

			db = &models.DB{
				StreamsMap: map[int]*models.Stream{
					1234: &models.Stream{
						Pid: 1234,
						EntitiesMap: map[uint64]*models.Entity{
							1: &models.Entity{ID: 1, Name: "FooBar", Index: 2},
							2: &models.Entity{ID: 2, Name: "Baah", Index: 1},
						},
					},
					5678: &models.Stream{Pid: 5678},
				},
				StreamKeys:        []int{1234, 5678},
				StreamEventSource: fakeStreamEventSource,
				EntityEventSource: fakeEntityEventSource,
			}
		})

		Describe("Streams", func() {
			It("returns all the streams found in the database", func() {
				Expect(db.Streams()).To(Equal([]models.Stream{
					*db.StreamsMap[1234],
					*db.StreamsMap[5678],
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
				Expect(db.Entity(1234, 1)).To(Equal(models.Entity{ID: 1, Name: "FooBar", Index: 2}))
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

		Describe("StreamEvent", func() {
			var eventsChannel chan models.StreamEvent

			BeforeEach(func() {
				eventsChannel = make(chan models.StreamEvent, 1)
				fakeStreamEventSource.SubscribeReturns(eventsChannel, 1234)
			})

			AfterEach(func() {
				close(eventsChannel)
			})

			It("returns a channel on which clients can receive events", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				payload := models.StreamEvent{
					StreamID: 123,
				}
				eventsChannel <- payload
				ch, err := db.StreamEvent(ctx)
				Expect(err).ToNot(HaveOccurred())

				var receivedPayload models.StreamEvent
				Expect(ch).To(Receive(&receivedPayload))
				Expect(receivedPayload).To(Equal(payload))
			})

			It("unsubscribes from the source when the context is done", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				_, err := db.StreamEvent(ctx)
				Expect(err).ToNot(HaveOccurred())

				cancel()
				Eventually(fakeStreamEventSource.UnsubscribeCallCount).Should(Equal(1))
				Expect(fakeStreamEventSource.UnsubscribeArgsForCall(0)).To(Equal(uint64(1234)))
			})
		})

		Describe("EntityEvent", func() {
			var eventsChannel chan models.EntityEvent

			BeforeEach(func() {
				eventsChannel = make(chan models.EntityEvent, 1)
				fakeEntityEventSource.SubscribeReturns(eventsChannel, 1234)
			})

			AfterEach(func() {
				close(eventsChannel)
			})

			It("returns a channel on which clients can receive events", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				payload := models.EntityEvent{
					StreamID: 123,
				}
				eventsChannel <- payload
				ch, err := db.EntityEvent(ctx)
				Expect(err).ToNot(HaveOccurred())

				var receivedPayload models.EntityEvent
				Expect(ch).To(Receive(&receivedPayload))
				Expect(receivedPayload).To(Equal(payload))
			})

			It("unsubscribes from the source when the context is done", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				_, err := db.EntityEvent(ctx)
				Expect(err).ToNot(HaveOccurred())

				cancel()
				Eventually(fakeEntityEventSource.UnsubscribeCallCount).Should(Equal(1))
				Expect(fakeEntityEventSource.UnsubscribeArgsForCall(0)).To(Equal(uint64(1234)))
			})
		})
	})

	Describe("Stream", func() {
		var stream *models.Stream

		BeforeEach(func() {
			stream = &models.Stream{
				Pid: 1234,
				EntitiesMap: map[uint64]*models.Entity{
					1: &models.Entity{ID: 1, Name: "FooBar", Index: 2},
					2: &models.Entity{ID: 2, Name: "Baah", Index: 1},
				},
			}
		})

		Describe("Entity", func() {
			It("returns the requested entity from the stream", func() {
				Expect(stream.Entity(1)).To(Equal(*stream.EntitiesMap[1]))
			})

			It("returns an error if the requested entity does not exist", func() {
				_, err := stream.Entity(3)
				Expect(err).To(MatchError("entity ID 3 not found"))
			})
		})

		Describe("Entities", func() {
			It("returns all entities found on the stream, sorted in order by index", func() {
				Expect(stream.Entities()).To(Equal([]models.Entity{
					*stream.EntitiesMap[2],
					*stream.EntitiesMap[1],
				}))
			})
		})
	})

	Describe("Timestamp", func() {
		It("marshals the provided time to the time since the Unix epoch in milliseconds", func() {
			t := time.Unix(101, 302000000)
			m := models.MarshalTimestamp(t)
			b := new(bytes.Buffer)
			m.MarshalGQL(b)
			Expect(b.String()).To(Equal("101302"))
		})
	})

	Describe("Uint", func() {
		It("marshals the provided uint to string", func() {
			m := models.MarshalUint(123456789000000)
			b := new(bytes.Buffer)
			m.MarshalGQL(b)
			Expect(b.String()).To(Equal("123456789000000"))
		})

		It("unmarshals the string to the expected uint", func() {
			u, err := models.UnmarshalUint("123456789000000")
			Expect(err).ToNot(HaveOccurred())
			Expect(u).To(Equal(uint64(123456789000000)))
		})

		It("unmarshals the JSON number to the expected uint", func() {
			u, err := models.UnmarshalUint(json.Number("123456789000000"))
			Expect(err).ToNot(HaveOccurred())
			Expect(u).To(Equal(uint64(123456789000000)))
		})

		It("unmarshals the integers to the expected uint", func() {
			u, err := models.UnmarshalUint(123456789)
			Expect(err).ToNot(HaveOccurred())
			Expect(u).To(Equal(uint64(123456789)))
		})

		It("errors if the data is not an integer type", func() {
			_, err := models.UnmarshalUint(1.2)
			Expect(err).To(MatchError(MatchRegexp(`.* is not a supported integer type`)))
		})
	})
})
