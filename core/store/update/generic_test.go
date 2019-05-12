package update_test

import (
	"time"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testVars struct {
	b       *xivnet.Block
	streams *store.Streams
	d       *datasheet.Collection

	streamID  int
	subjectID uint64
	entity    *models.Entity

	generator update.Generator
}

func genericSetup() (testEnv testVars) {
	testEnv.streamID = 1234
	testEnv.subjectID = 0x12345678

	testEnv.entity = &models.Entity{}

	testEnv.streams = &store.Streams{
		Map: map[int]*models.Stream{
			testEnv.streamID: &models.Stream{
				ID:          testEnv.streamID,
				CharacterID: testEnv.subjectID,
				EntitiesMap: map[uint64]*models.Entity{
					testEnv.subjectID: testEnv.entity,
					0x23456789:        nil,
				},
			},
		},
	}

	testEnv.d = new(datasheet.Collection)

	testEnv.generator = update.NewGenerator(testEnv.d)

	testEnv.b = &xivnet.Block{
		Length:    1234,
		SubjectID: uint32(testEnv.subjectID),
		CurrentID: 0x9ABCDEF0,
		IPCHeader: xivnet.IPCHeader{
			Opcode: 1234,
			Time:   time.Unix(102, 0),
		},
	}

	return
}

func entityValidationTests(testEnv *testVars, isEgress bool) {
	It("errors when the stream doesn't exist", func() {
		generator := testEnv.generator
		b := testEnv.b
		streams := testEnv.streams

		u := generator.Generate(1000, isEgress, b)
		Expect(u).ToNot(BeNil())

		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).To(MatchError(update.ErrorStreamNotFound))
		Expect(streamEvents).To(BeEmpty())
		Expect(entityEvents).To(BeEmpty())
	})

	It("errors when the entity doesn't exist", func() {
		generator := testEnv.generator
		b := testEnv.b
		streams := testEnv.streams
		streamID := testEnv.streamID

		b.SubjectID = 0x9ABCDEF0

		u := generator.Generate(streamID, isEgress, b)
		Expect(u).ToNot(BeNil())

		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).To(MatchError(update.ErrorEntityNotFound))
		Expect(streamEvents).To(BeEmpty())
		Expect(entityEvents).To(BeEmpty())
	})

	It("does nothing if the entity is nil", func() {
		generator := testEnv.generator
		b := testEnv.b
		streams := testEnv.streams
		streamID := testEnv.streamID

		b.SubjectID = 0x23456789

		u := generator.Generate(streamID, isEgress, b)
		Expect(u).ToNot(BeNil())

		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).To(BeNil())
		Expect(streamEvents).To(BeEmpty())
		Expect(entityEvents).To(BeEmpty())
	})

	It("does nothing if the stream's CharacterID is 0", func() {
		generator := testEnv.generator
		b := testEnv.b
		streams := testEnv.streams
		streamID := testEnv.streamID

		streams.Map[streamID].CharacterID = 0

		u := generator.Generate(streamID, isEgress, b)
		Expect(u).ToNot(BeNil())

		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).To(BeNil())
		Expect(streamEvents).To(BeEmpty())
		Expect(entityEvents).To(BeEmpty())
	})
}
