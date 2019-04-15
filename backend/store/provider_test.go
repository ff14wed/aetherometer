package store_test

import (
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/sibyl/backend/testhelpers"
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

type testUpdate func(*store.Streams) ([]models.StreamEvent, []models.EntityEvent, error)

func (t testUpdate) ModifyStore(s *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	return t(s)
}

var _ = Describe("Provider", func() {
	var (
		provider *store.Provider
		stream1  models.Stream
		stream2  models.Stream

		logBuf *testhelpers.LogBuffer
		once   sync.Once

		supervisor *suture.Supervisor
	)

	BeforeEach(func() {
		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err := zap.RegisterSink("providertest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
		logBuf.Reset()
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"providertest://"}
		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		stream1 = models.Stream{
			ID: 1234,
			EntitiesMap: map[uint64]*models.Entity{
				1: &models.Entity{ID: 1, Name: "FooBar", Index: 2},
				2: &models.Entity{ID: 2, Name: "Baah", Index: 1},
				3: nil,
			},
		}
		stream2 = models.Stream{ID: 5678}
		provider = store.NewProvider(
			logger,
			store.WithQueryTimeout(10*time.Millisecond),
			store.WithUpdateBufferSize(10),
			store.WithEventBufferSize(10),
			store.WithRequestBufferSize(10),
		)

		supervisor = suture.New("test-provider", suture.Spec{
			Log: func(line string) {
				_, _ = GinkgoWriter.Write([]byte(line))
			},
			FailureThreshold: 1,
		})
		supervisor.ServeBackground()
		_ = supervisor.Add(provider)

		provider.UpdatesChan() <- testUpdate(func(s *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
			s.Map[1234] = &stream1
			s.Map[5678] = &stream2
			s.KeyOrder = []int{5678, 1234}

			return nil, nil, nil
		})

		Eventually(provider.Streams).Should(HaveLen(2))
	})

	AfterEach(func() {
		supervisor.Stop()
	})

	It(`logs "Running" on startup`, func() {
		Eventually(logBuf).Should(gbytes.Say("store-provider.*Running"))
	})

	It(`logs "Stopping..." on shutdown`, func() {
		supervisor.Stop()
		Eventually(logBuf).Should(gbytes.Say("store-provider.*Stopping..."))
	})

	Describe("Streams", func() {
		It("returns all the streams found in the store", func() {
			Expect(provider.Streams()).To(Equal([]models.Stream{stream2, stream1}))
		})

		It("times out requests that take too long", func() {
			blockCh := make(chan struct{})
			provider.UpdatesChan() <- testUpdate(func(s *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
				By("blocking the provider service loop")
				<-blockCh
				<-blockCh
				return nil, nil, nil
			})

			Eventually(blockCh).Should(BeSent(struct{}{}))
			_, err := provider.Streams()
			Expect(err).To(MatchError(store.ErrRequestTimedOut))
			close(blockCh)
		})
	})

	Describe("Stream", func() {
		It("returns the requested stream from the store", func() {
			Expect(provider.Stream(5678)).To(Equal(stream2))
		})

		It("returns an error if the requested stream does not exist", func() {
			_, err := provider.Stream(2345)
			Expect(err).To(MatchError("stream ID 2345 not found"))
		})

		It("times out requests that take too long", func() {
			blockCh := make(chan struct{})
			provider.UpdatesChan() <- testUpdate(func(s *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
				By("blocking the provider service loop")
				<-blockCh
				<-blockCh
				return nil, nil, nil
			})

			Eventually(blockCh).Should(BeSent(struct{}{}))
			_, err := provider.Stream(2345)
			Expect(err).To(MatchError(store.ErrRequestTimedOut))
			close(blockCh)
		})
	})

	Describe("Entity", func() {
		It("returns the requested entity from the store", func() {
			Expect(provider.Entity(1234, 1)).To(Equal(models.Entity{ID: 1, Name: "FooBar", Index: 2}))
		})

		It("returns an error if the requested stream does not exist", func() {
			_, err := provider.Entity(2345, 1)
			Expect(err).To(MatchError("entity ID 1 not found in stream 2345"))
		})

		It("returns an error if the requested entity does not exist", func() {
			_, err := provider.Entity(1234, 4)
			Expect(err).To(MatchError("entity ID 4 not found in stream 1234"))
		})

		It("returns an error if the requested entity is nil", func() {
			_, err := provider.Entity(1234, 3)
			Expect(err).To(MatchError("entity ID 3 not found in stream 1234"))
		})

		It("times out requests that take too long", func() {
			blockCh := make(chan struct{})
			provider.UpdatesChan() <- testUpdate(func(s *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
				By("blocking the provider service loop")
				<-blockCh
				<-blockCh
				return nil, nil, nil
			})

			Eventually(blockCh).Should(BeSent(struct{}{}))
			_, err := provider.Entity(1234, 3)
			Expect(err).To(MatchError(store.ErrRequestTimedOut))
			close(blockCh)
		})
	})

	Describe("UpdatesChan", func() {
		It("consumes updates and applies them to the internal store", func() {
			provider.UpdatesChan() <- testUpdate(func(s *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
				s.Map[5678].CharacterID = 2345
				return nil, nil, nil
			})
			Eventually(func() models.Stream {
				s, _ := provider.Stream(5678)
				return s
			}).Should(Equal(models.Stream{ID: 5678, CharacterID: 2345}))
		})

		It("ignores nil updates", func() {
			provider.UpdatesChan() <- nil
			Consistently(func() models.Stream {
				s, _ := provider.Stream(5678)
				return s
			}).Should(Equal(models.Stream{ID: 5678}))
			Expect(logBuf).To(gbytes.Say("Running"))
			Consistently(logBuf).ShouldNot(gbytes.Say("store-provider"))
		})

		It("updates applied should not affect the result of already returned queries", func() {
			queriedStream, err := provider.Stream(5678)
			Expect(err).ToNot(HaveOccurred())

			provider.UpdatesChan() <- testUpdate(func(s *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
				s.Map[5678].CharacterID = 2345
				return nil, nil, nil
			})
			Eventually(func() models.Stream {
				s, _ := provider.Stream(5678)
				return s
			}).Should(Equal(models.Stream{ID: 5678, CharacterID: 2345}))

			Expect(queriedStream).To(Equal(models.Stream{ID: 5678}))

		})

		It("logs errors returned by the update", func() {
			provider.UpdatesChan() <- testUpdate(func(s *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
				return nil, nil, errors.New("kaboom")
			})
			Eventually(logBuf).Should(gbytes.Say(`ERROR.*store-provider.*Error applying update.*update.*testUpdate.*error.*kaboom`))
		})

		It("broadcasts stream events and entity events returned by the update", func() {
			streamEvents1, sub1 := provider.StreamEventSource().Subscribe()
			streamEvents2, sub2 := provider.StreamEventSource().Subscribe()

			entityEvents, sub3 := provider.EntityEventSource().Subscribe()

			defer func() {
				provider.StreamEventSource().Unsubscribe(sub1)
				provider.StreamEventSource().Unsubscribe(sub2)

				provider.EntityEventSource().Unsubscribe(sub3)
			}()

			provider.UpdatesChan() <- testUpdate(func(s *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
				return []models.StreamEvent{
						models.StreamEvent{StreamID: 1234},
						models.StreamEvent{StreamID: 5678},
					}, []models.EntityEvent{
						models.EntityEvent{StreamID: 1234, EntityID: 2},
						models.EntityEvent{StreamID: 5678, EntityID: 2},
					}, nil
			})

			Eventually(streamEvents1).Should(Receive(Equal(models.StreamEvent{StreamID: 1234})))
			Eventually(streamEvents1).Should(Receive(Equal(models.StreamEvent{StreamID: 5678})))
			Eventually(streamEvents2).Should(Receive(Equal(models.StreamEvent{StreamID: 1234})))
			Eventually(streamEvents2).Should(Receive(Equal(models.StreamEvent{StreamID: 5678})))

			Eventually(entityEvents).Should(Receive(Equal(models.EntityEvent{StreamID: 1234, EntityID: 2})))
			Eventually(entityEvents).Should(Receive(Equal(models.EntityEvent{StreamID: 5678, EntityID: 2})))
		})

		It("broadcasts events even after an error", func() {
			entityEvents, sub := provider.EntityEventSource().Subscribe()

			defer func() {
				provider.EntityEventSource().Unsubscribe(sub)
			}()

			provider.UpdatesChan() <- testUpdate(func(s *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
				return nil, []models.EntityEvent{
					models.EntityEvent{StreamID: 1234, EntityID: 2},
				}, errors.New("kaboom")
			})
			Eventually(entityEvents).Should(Receive(Equal(models.EntityEvent{StreamID: 1234, EntityID: 2})))
		})
	})
})
