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

var _ = Describe("RemoveEntity Update", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streams   *store.Streams
		streamID  int
		generator update.Generator
	)

	const removableID uint64 = 0x99999999

	BeforeEach(func() {
		*testEnv = genericSetup()
		b = testEnv.b
		streams = testEnv.streams
		streamID = testEnv.streamID
		generator = testEnv.generator

		removeEntityData := &datatypes.RemoveEntity{
			ID: uint32(removableID),
		}
		b.Data = removeEntityData
	})

	It("generates an update that removes the entity with the specified ID", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())
		Expect(streamEvents).To(BeEmpty())

		Expect(entityEvents).To(HaveLen(1))
		Expect(entityEvents[0].StreamID).To(Equal(streamID))
		Expect(entityEvents[0].EntityID).To(Equal(removableID))
		eventType, assignable := entityEvents[0].Type.(models.RemoveEntity)
		Expect(assignable).To(BeTrue())
		Expect(eventType.ID).To(Equal(removableID))

		Expect(streams.Map[streamID].EntitiesMap).To(HaveKey(removableID))
		Expect(streams.Map[streamID].EntitiesMap[removableID]).To(BeNil())
	})

	Context(`when the specified entity doesn't "exist"`, func() {
		const nonexistentID uint64 = 0x88888888
		BeforeEach(func() {
			removeEntityData := &datatypes.RemoveEntity{
				ID: uint32(nonexistentID),
			}
			b.Data = removeEntityData
		})

		It(`successfully "removes" the entity`, func() {
			Expect(streams.Map[streamID].EntitiesMap).ToNot(HaveKey(nonexistentID))

			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

			Expect(entityEvents).To(HaveLen(1))
			Expect(entityEvents[0].StreamID).To(Equal(streamID))
			Expect(entityEvents[0].EntityID).To(Equal(nonexistentID))
			eventType, assignable := entityEvents[0].Type.(models.RemoveEntity)
			Expect(assignable).To(BeTrue())
			Expect(eventType.ID).To(Equal(nonexistentID))

			Expect(streams.Map[streamID].EntitiesMap).To(HaveKey(nonexistentID))
			Expect(streams.Map[streamID].EntitiesMap[nonexistentID]).To(BeNil())

		})
	})

	streamValidationTests(testEnv, false)
})
