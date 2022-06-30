package update_test

import (
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/dealancer/validate.v2"
)

var _ = Describe("PlayerStats Update", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streams   *store.Streams
		streamID  int
		generator update.Generator

		expectedStats *models.Stats
	)

	BeforeEach(func() {
		*testEnv = genericSetup()
		b = testEnv.b
		streams = testEnv.streams
		streamID = testEnv.streamID
		generator = testEnv.generator

		expectedStats = &models.Stats{
			Strength: 1000,
		}

		updatePlayerStatsData := &datatypes.PlayerStats{
			Strength: 1000,
		}
		b.Data = updatePlayerStatsData
	})

	It("generates an update that sets the stats", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())
		Expect(streamEvents).To(HaveLen(1))
		Expect(streamEvents[0].StreamID).To(Equal(streamID))
		eventType, assignable := streamEvents[0].Type.(models.UpdateStats)
		Expect(assignable).To(BeTrue())
		Expect(eventType.Stats).To(Equal(expectedStats))

		Expect(streams.Map[streamID].Stats).To(Equal(expectedStats))

		Expect(entityEvents).To(BeEmpty())

		Expect(validate.Validate(streamEvents)).To(Succeed())
		Expect(validate.Validate(streams)).To(Succeed())
	})

	streamValidationTests(testEnv, false)
})
