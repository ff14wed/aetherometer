package update_test

import (
	"time"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/dealancer/validate.v2"
)

var _ = Describe("EffectResult Update", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streams   *store.Streams
		streamID  int
		subjectID uint64
		entity    *models.Entity
		generator update.Generator

		expectedResources *models.Resources
		expectedStatuses  []*models.Status
	)

	BeforeEach(func() {
		*testEnv = genericSetup()
		b = testEnv.b
		streams = testEnv.streams
		streamID = testEnv.streamID
		subjectID = testEnv.subjectID
		entity = testEnv.entity
		generator = testEnv.generator

		testEnv.d.StatusData = map[uint32]datasheet.Status{
			123: {Key: 123, Name: "Foo"},
			456: {Key: 456, Name: "Bar"},
			789: {Key: 789, Name: "Baz"},
		}

		testEnv.entity.Resources.MaxMp = 300
		expectedResources = &models.Resources{
			Hp:       100,
			MaxHp:    100,
			Mp:       200,
			MaxMp:    300,
			Tp:       0,
			LastTick: b.Time,
		}

		expectedStatuses = []*models.Status{
			{ID: 123, Name: "Foo", StartedTime: b.Time, Duration: time.Unix(0, 0), LastTick: b.Time},
			{ID: 456, Name: "Bar", StartedTime: b.Time, Duration: time.Unix(0, 0), LastTick: b.Time},
			{ID: 789, Name: "Baz", StartedTime: b.Time, Duration: time.Unix(0, 0), LastTick: b.Time},
		}

		effectResultData := &datatypes.EffectResult{
			ActorID: uint32(testEnv.subjectID),

			CurrentHP: 100,
			MaxHP:     100,
			CurrentMP: 200,

			Count: 3,

			Entries: [4]datatypes.EffectResultEntry{
				{Index: 0, EffectID: 123},
				{Index: 1, EffectID: 456},
				{Index: 2, EffectID: 789},
			},
		}
		b.Data = effectResultData
	})

	It("generates an update that sets the entity's resources", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())
		Expect(streamEvents).To(BeEmpty())

		Expect(entityEvents).ToNot(BeEmpty())
		Expect(entityEvents[0].StreamID).To(Equal(streamID))
		Expect(entityEvents[0].EntityID).To(Equal(subjectID))
		eventType, assignable := entityEvents[0].Type.(models.UpdateResources)
		Expect(assignable).To(BeTrue())
		Expect(eventType.Resources).To(Equal(expectedResources))

		Expect(entity.Resources).To(Equal(expectedResources))

		Expect(validate.Validate(entityEvents)).To(Succeed())
		Expect(validate.Validate(streams)).To(Succeed())
	})

	It("generates an update that updates the entity's statuses", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())
		Expect(streamEvents).To(BeEmpty())

		statusEvents := entityEvents[1:]
		Expect(statusEvents).To(HaveLen(3))

		receivedEvents := []models.UpsertStatus{}
		for _, e := range statusEvents {
			Expect(e.StreamID).To(Equal(streamID))
			Expect(e.EntityID).To(Equal(subjectID))
			eventType, assignable := e.Type.(models.UpsertStatus)
			Expect(assignable).To(BeTrue())
			receivedEvents = append(receivedEvents, eventType)
		}

		expectedEvents := []interface{}{}
		for i, e := range expectedStatuses {
			expectedEvents = append(expectedEvents, models.UpsertStatus{Index: i, Status: e})
		}
		Expect(receivedEvents).To(ConsistOf(expectedEvents...))

		Expect(entity.Statuses).To(Equal(expectedStatuses))

		Expect(validate.Validate(entityEvents)).To(Succeed())
		Expect(validate.Validate(streams)).To(Succeed())
	})

	entityValidationTests(testEnv, false)
})
