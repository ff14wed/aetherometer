package update_test

import (
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/sibyl/backend/store/update"
	"github.com/ff14wed/xivnet/v2"
	"github.com/ff14wed/xivnet/v2/datatypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EquipChange Update", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streams   *store.Streams
		streamID  int
		subjectID uint64
		entity    *models.Entity
		generator update.Generator

		expectedClass models.ClassJob
	)

	BeforeEach(func() {
		*testEnv = genericSetup()
		b = testEnv.b
		streams = testEnv.streams
		streamID = testEnv.streamID
		subjectID = testEnv.subjectID
		entity = testEnv.entity
		generator = testEnv.generator

		expectedClass = models.ClassJob{
			ID: 0x12,
		}

		equipChangeData := &datatypes.EquipChange{
			ClassJob: 0x12,
		}
		b.Data = equipChangeData
	})

	It("generates an update that sets the entity's ClassJob", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())
		Expect(streamEvents).To(BeEmpty())

		Expect(entityEvents).To(ConsistOf(models.EntityEvent{
			StreamID: streamID,
			EntityID: subjectID,
			Type: models.UpdateClass{
				ClassJob: expectedClass,
			},
		}))

		Expect(entity.ClassJob).To(Equal(expectedClass))
	})

	entityValidationTests(testEnv, false)
})
