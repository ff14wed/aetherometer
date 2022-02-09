package stream_test

import (
	"errors"

	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/stream"
	"github.com/ff14wed/aetherometer/core/stream/streamfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("BuildAdapterInventory", func() {
	var (
		hookBuilder, testBuilder *streamfakes.FakeAdapterBuilder
		inventory                []stream.AdapterInfo
		cfg                      config.Config

		streamUp   chan<- stream.Provider
		streamDown chan<- int
		logger     *zap.Logger
	)

	BeforeEach(func() {
		hookBuilder = new(streamfakes.FakeAdapterBuilder)
		testBuilder = new(streamfakes.FakeAdapterBuilder)
		inventory = []stream.AdapterInfo{
			{Name: "Hook", Builder: hookBuilder},
			{Name: "test", Builder: testBuilder},
		}

		streamUp = make(chan stream.Provider)
		streamDown = make(chan int)
		logger = zap.NewExample()
	})

	Context("when an adapter is not enabled in the configuration", func() {
		BeforeEach(func() {
			cfg = config.Config{
				Adapters: config.Adapters{
					Hook: config.HookConfig{Enabled: false},
				},
			}
		})

		It("skips loading the configuration for this adapter and loads the rest", func() {
			adapters, err := stream.BuildAdapterInventory(inventory, cfg, streamUp, streamDown, logger)
			Expect(err).ToNot(HaveOccurred())
			Expect(adapters).To(HaveLen(1))
			Expect(adapters).To(HaveKey("test"))
			Expect(testBuilder.BuildCallCount()).To(Equal(1))
		})
	})

	Context("when the adapter is enabled in the configuration", func() {
		BeforeEach(func() {
			cfg = config.Config{
				Adapters: config.Adapters{
					Hook: config.HookConfig{Enabled: true},
				},
			}
		})

		It("calls Build on all of the preconfigured AdapterInfos", func() {
			adapters, err := stream.BuildAdapterInventory(inventory, cfg, streamUp, streamDown, logger)
			Expect(err).ToNot(HaveOccurred())
			Expect(adapters).To(HaveLen(2))
			Expect(adapters).To(HaveKey("Hook"))
			Expect(adapters).To(HaveKey("test"))
			Expect(hookBuilder.BuildCallCount()).To(Equal(1))
			Expect(testBuilder.BuildCallCount()).To(Equal(1))

			var (
				streamUpArg   chan<- stream.Provider
				streamDownArg chan<- int
				loggerArg     *zap.Logger
			)
			streamUpArg, streamDownArg, loggerArg = hookBuilder.BuildArgsForCall(0)
			Expect(streamUpArg).To(Equal(streamUp))
			Expect(streamDownArg).To(Equal(streamDown))
			Expect(loggerArg).To(Equal(logger))

			streamUpArg, streamDownArg, loggerArg = testBuilder.BuildArgsForCall(0)
			Expect(streamUpArg).To(Equal(streamUp))
			Expect(streamDownArg).To(Equal(streamDown))
			Expect(loggerArg).To(Equal(logger))
		})

		Context("when loading configuration fails", func() {
			BeforeEach(func() {
				cfg = config.Config{
					Adapters: config.Adapters{
						Hook: config.HookConfig{Enabled: true},
					},
				}
				hookBuilder.LoadConfigReturns(errors.New("boom"))
			})

			It("does not call Build and errors without building the rest of the adapters", func() {
				adapters, err := stream.BuildAdapterInventory(inventory, cfg, streamUp, streamDown, logger)
				Expect(err).To(MatchError("error creating adapter Hook: boom"))
				Expect(adapters).To(BeEmpty())
				Expect(hookBuilder.BuildCallCount()).To(Equal(0))
				Expect(testBuilder.BuildCallCount()).To(Equal(0))
			})
		})
	})
})
