package models_test

import (
	"context"
	"errors"

	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/models/modelsfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Models", func() {
	Describe("Resolver", func() {
		var (
			fakeStreamEventSource *modelsfakes.FakeStreamEventSource
			fakeEntityEventSource *modelsfakes.FakeEntityEventSource
			fakeStoreProvider     *modelsfakes.FakeStoreProvider
			fakeAuthProvider      *modelsfakes.FakeAuthProvider
			resolver              *models.Resolver
			stream1, stream2      models.Stream
		)

		BeforeEach(func() {
			fakeStreamEventSource = new(modelsfakes.FakeStreamEventSource)
			fakeEntityEventSource = new(modelsfakes.FakeEntityEventSource)
			fakeStoreProvider = new(modelsfakes.FakeStoreProvider)
			fakeAuthProvider = new(modelsfakes.FakeAuthProvider)

			stream1 = models.Stream{
				ID: 1234,
				EntitiesMap: map[uint64]*models.Entity{
					1: {ID: 1, Name: "FooBar", Index: 2},
					2: {ID: 2, Name: "Baah", Index: 1},
				},
			}
			stream2 = models.Stream{ID: 5678}
			fakeStoreProvider.StreamsReturns([]models.Stream{stream1, stream2}, nil)

			fakeStoreProvider.StreamStub = func(streamID int) (*models.Stream, error) {
				if streamID == 1234 {
					return &stream1, nil
				} else if streamID == 5678 {
					return &stream2, nil
				}
				return nil, errors.New("not found")
			}

			fakeStoreProvider.EntityStub = func(streamID int, entityID uint64) (*models.Entity, error) {
				s, err := fakeStoreProvider.Stream(streamID)
				if err != nil {
					return nil, err
				}
				if e, found := s.EntitiesMap[entityID]; found {
					return e, nil
				}
				return nil, errors.New("not found")
			}

			fakeStoreProvider.StreamEventSourceReturns(fakeStreamEventSource)
			fakeStoreProvider.EntityEventSourceReturns(fakeEntityEventSource)

			resolver = models.NewResolver(fakeStoreProvider, fakeAuthProvider, nil)
		})

		Describe("Streams", func() {
			It("returns all the streams found in the store", func() {
				Expect(resolver.Query().Streams(context.Background())).To(Equal(
					[]models.Stream{stream1, stream2},
				))
			})

			Context("when the request is not authorized", func() {
				BeforeEach(func() {
					fakeAuthProvider.AuthorizePluginTokenReturns(errors.New("Boom"))
				})

				It("returns an authorization error", func() {
					streams, err := resolver.Query().Streams(context.Background())
					Expect(err).To(MatchError("Boom"))
					Expect(streams).To(BeZero())
				})
			})
		})

		Describe("Stream", func() {
			It("returns the requested stream from the store", func() {
				Expect(resolver.Query().Stream(context.Background(), 5678)).To(Equal(&stream2))
			})

			It("returns an error if the requested stream does not exist", func() {
				_, err := resolver.Query().Stream(context.Background(), 2345)
				Expect(err).To(MatchError("not found"))
			})

			Context("when the request is not authorized", func() {
				BeforeEach(func() {
					fakeAuthProvider.AuthorizePluginTokenReturns(errors.New("Boom"))
				})

				It("returns an authorization error", func() {
					s, err := resolver.Query().Stream(context.Background(), 2345)
					Expect(err).To(MatchError("Boom"))
					Expect(s).To(BeZero())
				})
			})
		})

		Describe("Entity", func() {
			It("returns the requested entity from the store", func() {
				Expect(resolver.Query().Entity(context.Background(), 1234, 1)).To(Equal(
					&models.Entity{ID: 1, Name: "FooBar", Index: 2},
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

			Context("when the request is not authorized", func() {
				BeforeEach(func() {
					fakeAuthProvider.AuthorizePluginTokenReturns(errors.New("Boom"))
				})

				It("returns an authorization error", func() {
					e, err := resolver.Query().Entity(context.Background(), 1234, 1)
					Expect(err).To(MatchError("Boom"))
					Expect(e).To(BeZero())
				})
			})
		})

		Describe("StreamEvent", func() {
			var eventsChannel chan *models.StreamEvent

			BeforeEach(func() {
				eventsChannel = make(chan *models.StreamEvent, 1)
				fakeStreamEventSource.SubscribeReturns(eventsChannel, 1234)
			})

			AfterEach(func() {
				close(eventsChannel)
			})

			It("returns a channel on which clients can receive events", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				payload := &models.StreamEvent{
					StreamID: 123,
				}
				eventsChannel <- payload
				ch, err := resolver.Subscription().StreamEvent(ctx)
				Expect(err).ToNot(HaveOccurred())

				var receivedPayload *models.StreamEvent
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

			Context("when the request is not authorized", func() {
				BeforeEach(func() {
					fakeAuthProvider.AuthorizePluginTokenReturns(errors.New("Boom"))
				})

				It("returns an authorization error", func() {
					ch, err := resolver.Subscription().StreamEvent(context.Background())
					Expect(err).To(MatchError("Boom"))
					Expect(ch).To(BeZero())
				})
			})
		})

		Describe("EntityEvent", func() {
			var eventsChannel chan *models.EntityEvent

			BeforeEach(func() {
				eventsChannel = make(chan *models.EntityEvent, 1)
				fakeEntityEventSource.SubscribeReturns(eventsChannel, 1234)
			})

			AfterEach(func() {
				close(eventsChannel)
			})

			It("returns a channel on which clients can receive events", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				payload := &models.EntityEvent{
					StreamID: 123,
				}
				eventsChannel <- payload
				ch, err := resolver.Subscription().EntityEvent(ctx)
				Expect(err).ToNot(HaveOccurred())

				var receivedPayload *models.EntityEvent
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

			Context("when the request is not authorized", func() {
				BeforeEach(func() {
					fakeAuthProvider.AuthorizePluginTokenReturns(errors.New("Boom"))
				})

				It("returns an authorization error", func() {
					ch, err := resolver.Subscription().EntityEvent(context.Background())
					Expect(err).To(MatchError("Boom"))
					Expect(ch).To(BeZero())
				})
			})
		})

		Describe("SendStreamRequest", func() {
			var (
				requestedStreamID int
				requestedData     []byte
			)

			Context("when the request handler exists", func() {
				BeforeEach(func() {
					requestedStreamID = 0
					requestedData = nil
					resolver = models.NewResolver(fakeStoreProvider, fakeAuthProvider, func(streamID int, data []byte) (string, error) {
						requestedStreamID = streamID
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

					Expect(requestedStreamID).To(Equal(123))
					Expect(requestedData).To(Equal([]byte("hello")))
				})

				Context("when the handler errors", func() {
					BeforeEach(func() {
						resolver = models.NewResolver(fakeStoreProvider, fakeAuthProvider, func(streamID int, data []byte) (string, error) {
							return "", errors.New("kaboom")
						})
					})

					It("returns the handler's error", func() {
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

				Context("when the request is not authorized", func() {
					BeforeEach(func() {
						fakeAuthProvider.AuthorizePluginTokenReturns(errors.New("Boom"))
					})

					It("returns an authorization error", func() {
						resp, err := resolver.Mutation().SendStreamRequest(
							context.Background(),
							models.StreamRequest{
								StreamID: 123,
								Data:     "hello",
							},
						)
						Expect(resp).To(BeEmpty())
						Expect(err).To(MatchError("Boom"))
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

		Describe("CreateAdminToken", func() {
			It("calls the auth provider", func() {
				ctx := context.Background()
				fakeAuthProvider.CreateAdminTokenReturns("some-token", nil)

				Expect(resolver.Mutation().CreateAdminToken(ctx)).To(Equal(
					"some-token",
				))

				Expect(fakeAuthProvider.CreateAdminTokenCallCount()).To(Equal(1))
				ctxArg := fakeAuthProvider.CreateAdminTokenArgsForCall(0)
				Expect(ctxArg).To(Equal(ctx))
			})
		})

		Describe("AddPlugin", func() {
			It("calls the auth provider", func() {
				ctx := context.Background()
				pluginURL := "some-plugin-url"
				fakeAuthProvider.AddPluginReturns("some-token", nil)

				Expect(resolver.Mutation().AddPlugin(ctx, pluginURL)).To(Equal(
					"some-token",
				))

				Expect(fakeAuthProvider.AddPluginCallCount()).To(Equal(1))
				ctxArg, pluginURLArg := fakeAuthProvider.AddPluginArgsForCall(0)
				Expect(ctxArg).To(Equal(ctx))
				Expect(pluginURLArg).To(Equal(pluginURL))
			})
		})

		Describe("RemovePlugin", func() {
			It("calls the auth provider", func() {
				ctx := context.Background()
				authToken := "some-auto-token"
				fakeAuthProvider.RemovePluginReturns(true, nil)

				Expect(resolver.Mutation().RemovePlugin(ctx, authToken)).To(BeTrue())

				Expect(fakeAuthProvider.RemovePluginCallCount()).To(Equal(1))
				ctxArg, authTokenArg := fakeAuthProvider.RemovePluginArgsForCall(0)
				Expect(ctxArg).To(Equal(ctx))
				Expect(authTokenArg).To(Equal(authToken))
			})
		})
	})
})
