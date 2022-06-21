package update_test

import (
	"time"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/dealancer/validate.v2"
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

	testEnv.entity = &models.Entity{
		ID: testEnv.subjectID, Name: "Test Subject",
		ClassJob:  &models.ClassJob{},
		Resources: &models.Resources{},
		Location:  &models.Location{},
		Statuses: []*models.Status{
			&models.Status{ID: 1},
		},
	}

	testEnv.streams = &store.Streams{
		Map: map[int]*models.Stream{
			testEnv.streamID: {
				ID:          testEnv.streamID,
				ServerID:    2000,
				InstanceNum: 1000,

				CharacterID:  testEnv.subjectID,
				CurrentWorld: models.World{ID: 123, Name: "Foo"},
				HomeWorld:    models.World{ID: 456, Name: "Bar"},

				EntitiesMap: map[uint64]*models.Entity{
					testEnv.subjectID: testEnv.entity,
					0x23456789:        nil,
					0x99999999: {
						ID: 0x99999999, Index: 123,
						ClassJob:  &models.ClassJob{},
						Resources: &models.Resources{},
						Location:  &models.Location{},
					},
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
			Opcode:   1234,
			ServerID: 5678,
			Time:     time.Unix(102, 0),
		},
	}

	return
}

func streamValidationTests(testEnv *testVars, isEgress bool) {
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

		Expect(validate.Validate(streams)).To(Succeed())
	})
}

func entityValidationTests(testEnv *testVars, isEgress bool) {
	streamValidationTests(testEnv, isEgress)

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

		Expect(validate.Validate(streams)).To(Succeed())
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

		Expect(validate.Validate(streams)).To(Succeed())
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

		Expect(validate.Validate(streams)).To(Succeed())
	})
}
