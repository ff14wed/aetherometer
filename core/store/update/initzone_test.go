package update_test

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InitZone Update", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streams   *store.Streams
		d         *datasheet.Collection
		streamID  int
		generator update.Generator

		expectedPlace models.Place
	)

	BeforeEach(func() {
		*testEnv = genericSetup()
		b = testEnv.b
		streams = testEnv.streams
		d = testEnv.d
		streamID = testEnv.streamID
		generator = testEnv.generator

		d.MapData = datasheet.MapStore{
			Maps: map[uint16]datasheet.MapInfo{
				14: {
					Key: 14, ID: "w1t2/01", SizeFactor: 200, PlaceName: 41,
					PlaceNameSub: 373, TerritoryType: 131,
				},
				73: {
					Key: 73, ID: "w1t2/02", SizeFactor: 200, PlaceName: 41,
					PlaceNameSub: 698, TerritoryType: 131,
				},
			},
			Territories: map[uint16]datasheet.TerritoryInfo{
				131: {Key: 131, Name: "w1t2", Map: 14},
			},
			PlaceNames: map[uint16]datasheet.PlaceName{
				41:  {Key: 41, Name: "Ul'dah - Steps of Thal"},
				373: {Key: 373, Name: "Merchant Strip"},
				698: {Key: 698, Name: "Hustings Strip"},
			},
		}

		expectedPlace = models.Place{
			MapID:       14,
			TerritoryID: 131,
			Maps: []models.MapInfo{
				{
					Key:           14,
					ID:            "w1t2/01",
					SizeFactor:    200,
					PlaceName:     "Ul'dah - Steps of Thal",
					PlaceNameSub:  "Merchant Strip",
					TerritoryType: "w1t2",
				},
				{
					Key:           73,
					ID:            "w1t2/02",
					SizeFactor:    200,
					PlaceName:     "Ul'dah - Steps of Thal",
					PlaceNameSub:  "Hustings Strip",
					TerritoryType: "w1t2",
				},
			},
		}

		initZoneData := &datatypes.InitZone{
			TerritoryTypeID: 131,
		}
		b.Data = initZoneData
	})

	It("generates an update that sets the current server ID, character ID, and the instance num", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		streamEvents, _, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())

		Expect(streamEvents).To(HaveLen(2))
		Expect(streamEvents[0].StreamID).To(Equal(streamID))
		eventType, assignable := streamEvents[0].Type.(models.UpdateIDs)
		Expect(assignable).To(BeTrue())
		Expect(eventType.ServerID).To(Equal(int(b.ServerID)))
		Expect(eventType.CharacterID).To(Equal(uint64(b.CurrentID)))
		Expect(eventType.InstanceNum).To(Equal(0))

		Expect(streams.Map[streamID].ServerID).To(Equal(int(b.ServerID)))
		Expect(streams.Map[streamID].CharacterID).To(Equal(uint64(b.CurrentID)))
		Expect(streams.Map[streamID].InstanceNum).To(Equal(0))
	})

	It("generates an update that changes the place", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		streamEvents, _, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())

		Expect(streamEvents).To(HaveLen(2))
		Expect(streamEvents[1].StreamID).To(Equal(streamID))
		eventType, assignable := streamEvents[1].Type.(models.UpdateMap)
		Expect(assignable).To(BeTrue())
		Expect(eventType.Place).To(Equal(expectedPlace))

		Expect(streams.Map[streamID].Place).To(Equal(expectedPlace))
	})

	It("generates an update that clears the entity map", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		_, entityEvents, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())

		Expect(entityEvents).To(HaveLen(1))
		Expect(entityEvents[0].StreamID).To(Equal(streamID))
		Expect(entityEvents[0].EntityID).To(BeZero())
		eventType, assignable := entityEvents[0].Type.(models.SetEntities)
		Expect(assignable).To(BeTrue())
		Expect(eventType.Entities).To(BeNil())

		Expect(streams.Map[streamID].EntitiesMap).To(BeEmpty())
	})

	streamValidationTests(testEnv, false)
})
