package update_test

import (
	"time"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/sibyl/backend/store/update"
	"github.com/ff14wed/xivnet"
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
				PID: testEnv.streamID,
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
		Length: 1234,
		Header: xivnet.BlockHeader{
			SubjectID: uint32(testEnv.subjectID),
			CurrentID: 0x9ABCDEF0,
			Opcode:    1234,
			Time:      time.Unix(12, 0),
		},
	}
	return
}

func entityValidationTests(testEnv *testVars) {
	It("errors when the stream doesn't exist", func() {
		generator := testEnv.generator
		b := testEnv.b
		streams := testEnv.streams

		u := generator.Generate(1000, false, b)
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

		b.Header.SubjectID = 0x9ABCDEF0

		u := generator.Generate(streamID, false, b)
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

		b.Header.SubjectID = 0x23456789

		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())

		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).To(BeNil())
		Expect(streamEvents).To(BeEmpty())
		Expect(entityEvents).To(BeEmpty())
	})
}
