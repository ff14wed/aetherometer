package models_test

import (
	"context"
	"errors"

	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/models/modelsfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Models", func() {
	Describe("Resolver", func() {
		var (
			fakeStreamEventSource *modelsfakes.FakeStreamEventSource
			fakeEntityEventSource *modelsfakes.FakeEntityEventSource
			fakeStoreProvider     *modelsfakes.FakeStoreProvider
			resolver              *models.Resolver
			stream1, stream2      models.Stream
		)

		BeforeEach(func() {
			fakeStreamEventSource = new(modelsfakes.FakeStreamEventSource)
			fakeEntityEventSource = new(modelsfakes.FakeEntityEventSource)
			fakeStoreProvider = new(modelsfakes.FakeStoreProvider)

			stream1 = models.Stream{
				Pid: 1234,
				EntitiesMap: map[uint64]*models.Entity{
					1: &models.Entity{ID: 1, Name: "FooBar", Index: 2},
					2: &models.Entity{ID: 2, Name: "Baah", Index: 1},
				},
			}
			stream2 = models.Stream{Pid: 5678}
			fakeStoreProvider.StreamsReturns([]models.Stream{stream1, stream2}, nil)

			fakeStoreProvider.StreamStub = func(streamID int) (models.Stream, error) {
				if streamID == 1234 {
					return stream1, nil
				} else if streamID == 5678 {
					return stream2, nil
				}
				return models.Stream{}, errors.New("not found")
			}

			fakeStoreProvider.EntityStub = func(streamID int, entityID uint64) (models.Entity, error) {
				s, err := fakeStoreProvider.Stream(streamID)
				if err != nil {
					return models.Entity{}, err
				}
				if e, found := s.EntitiesMap[entityID]; found {
					return *e, nil
				}
				return models.Entity{}, errors.New("not found")
			}

			fakeStoreProvider.StreamEventSourceReturns(fakeStreamEventSource)
			fakeStoreProvider.EntityEventSourceReturns(fakeEntityEventSource)

			resolver = models.NewResolver(fakeStoreProvider, nil)
		})

		Describe("Streams", func() {
			It("returns all the streams found in the database", func() {
				Expect(resolver.Query().Streams(context.Background())).To(Equal(
					[]models.Stream{stream1, stream2},
				))
			})
		})

		Describe("Stream", func() {
			It("returns the requested stream from the database", func() {
				Expect(resolver.Query().Stream(context.Background(), 5678)).To(Equal(stream2))
			})

			It("returns an error if the requested stream does not exist", func() {
				_, err := resolver.Query().Stream(context.Background(), 2345)
				Expect(err).To(MatchError("not found"))
			})
		})

		Describe("Entity", func() {
			It("returns the requested entity from the database", func() {
				Expect(resolver.Query().Entity(context.Background(), 1234, 1)).To(Equal(
					models.Entity{ID: 1, Name: "FooBar", Index: 2},
				))
			})

			It("returns an error if the requested stream does not exist", func() {
				_, err := resolver.Query().Entity(context.Background(), 2345, 1)
				Expect(err).To(MatchError("not found"))
			})

			It("returns an error if the requested entity does not exist", func() {
				_, err := resolver.Query().Entity(context.Background(), 1234, 3)
				Expect(err).To(MatchError("not found"))
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
				ch, err := resolver.Subscription().StreamEvent(ctx)
				Expect(err).ToNot(HaveOccurred())

				var receivedPayload models.StreamEvent
				Expect(ch).To(Receive(&receivedPayload))
				Expect(receivedPayload).To(Equal(payload))
			})

			It("unsubscribes from the source when the context is done", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				_, err := resolver.Subscription().StreamEvent(ctx)
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
				ch, err := resolver.Subscription().EntityEvent(ctx)
				Expect(err).ToNot(HaveOccurred())

				var receivedPayload models.EntityEvent
				Expect(ch).To(Receive(&receivedPayload))
				Expect(receivedPayload).To(Equal(payload))
			})

			It("unsubscribes from the source when the context is done", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				_, err := resolver.Subscription().EntityEvent(ctx)
				Expect(err).ToNot(HaveOccurred())

				cancel()
				Eventually(fakeEntityEventSource.UnsubscribeCallCount).Should(Equal(1))
				Expect(fakeEntityEventSource.UnsubscribeArgsForCall(0)).To(Equal(uint64(1234)))
			})
		})

		Describe("SendStreamRequest", func() {
			var (
				requestedPid  int
				requestedData []byte
			)

			Context("when the request handler exists", func() {
				BeforeEach(func() {
					requestedPid = 0
					requestedData = nil
					resolver = models.NewResolver(fakeStoreProvider, func(pid int, data []byte) (string, error) {
						requestedPid = pid
						requestedData = data
						return "Success", nil
					})
				})

				It("successfully calls the provided handler", func() {
					resp, err := resolver.Mutation().SendStreamRequest(
						context.Background(),
						models.StreamRequest{
							StreamID: 123,
							Data:     "hello",
						},
					)
					Expect(resp).To(Equal("Success"))
					Expect(err).To(BeNil())

					Expect(requestedPid).To(Equal(123))
					Expect(requestedData).To(Equal([]byte("hello")))
				})

				Context("when the handler errors", func() {
					BeforeEach(func() {
						resolver = models.NewResolver(fakeStoreProvider, func(pid int, data []byte) (string, error) {
							return "", errors.New("kaboom")
						})
					})

					It("successfully calls the provided handler", func() {
						resp, err := resolver.Mutation().SendStreamRequest(
							context.Background(),
							models.StreamRequest{
								StreamID: 123,
								Data:     "hello",
							})
						Expect(resp).To(BeEmpty())
						Expect(err).To(MatchError("kaboom"))
					})
				})
			})

			Context("when the request handler is missing", func() {
				It("returns an error", func() {
					resp, err := resolver.Mutation().SendStreamRequest(
						context.Background(),
						models.StreamRequest{
							StreamID: 123,
							Data:     "hello",
						})
					Expect(resp).To(BeEmpty())
					Expect(err).To(MatchError("Request handler is missing"))
				})
			})
		})
	})
})
