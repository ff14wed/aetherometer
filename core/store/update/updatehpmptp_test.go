package update_test

import (
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateHPMPTP Update", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streams   *store.Streams
		streamID  int
		subjectID uint64
		entity    *models.Entity
		generator update.Generator

		expectedResources models.Resources
	)

	BeforeEach(func() {
		*testEnv = genericSetup()
		b = testEnv.b
		streams = testEnv.streams
		streamID = testEnv.streamID
		subjectID = testEnv.subjectID
		entity = testEnv.entity
		generator = testEnv.generator

		expectedResources = models.Resources{
			Hp:       100,
			Mp:       200,
			Tp:       300,
			LastTick: b.Time,
		}

		updateHPMPTPData := &datatypes.UpdateHPMPTP{
			HP: 100,
			MP: 200,
			TP: 300,
		}
		b.Data = updateHPMPTPData
	})

	It("generates an update that sets the entity's resources", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())
		Expect(streamEvents).To(BeEmpty())

		Expect(entityEvents).To(HaveLen(1))
		Expect(entityEvents[0].StreamID).To(Equal(streamID))
		Expect(entityEvents[0].EntityID).To(Equal(subjectID))
		eventType, assignable := entityEvents[0].Type.(models.UpdateResources)
		Expect(assignable).To(BeTrue())
		Expect(eventType.Resources).To(Equal(expectedResources))

		Expect(entity.Resources).To(Equal(expectedResources))
	})

	entityValidationTests(testEnv, false)
})
