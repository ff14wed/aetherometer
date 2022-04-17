package update_test

import (
	"encoding/json"
	"math"
	"time"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/aetherometer/core/testassets"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	"gopkg.in/dealancer/validate.v2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
)

var _ = Describe("Spawn Update", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streams   *store.Streams
		d         *datasheet.Collection
		streamID  int
		subjectID uint64
		generator update.Generator

		playerSpawnData *datatypes.PlayerSpawn

		expectedEntityFields gstruct.Fields
	)

	const removableID uint64 = 0x99999999

	BeforeEach(func() {
		*testEnv = genericSetup()
		b = testEnv.b
		streams = testEnv.streams
		d = testEnv.d
		streamID = testEnv.streamID
		subjectID = testEnv.subjectID
		generator = testEnv.generator

		playerSpawnData = &datatypes.PlayerSpawn{
			Title: 0x1234, U1b: 0x1234, CurrentWorld: 123, HomeWorld: 456,
			GMRank: 0x12, U3c: 0x9A, U4: 0x12, OnlineStatus: 0x12,
			Pose: 0x12, U5a: 0x12, U5b: 0x34, U5c: 0x12,
			TargetID: 0x9ABC, U6: 0x1234, U7: 0x1234,

			WeaponMain: datatypes.WeaponGear{
				Model1: 1, Model2: 2, Model3: 3, Model4: 4,
			},
			WeaponSub: datatypes.WeaponGear{
				Model1: 1, Model2: 2, Model3: 3, Model4: 4,
			},
			CraftSub: datatypes.WeaponGear{
				Model1: 1, Model2: 2, Model3: 3, Model4: 4,
			},
			BNPCBase: 0x5678, BNPCName: 0x5678,
			U18: 0x1234, U19: 0x1234, DirectorID: 0x1234,
			OwnerID: 0x9ABC,
			UnkID3:  0xDEF0,

			CurrentHP: 29000, DisplayFlags: 256, FateID: 0x1234, MaxHP: 30000,
			CurrentMP: 11000, MaxMP: 12000,

			ModelChara: 0x5678, Direction: 0x7FFF, MountID: 0,
			Minion: 0x1234, Index: 10, State: 1, Emote: 0x12, Type: 1,
			Subtype: 4, Voice: 0x12,

			EnemyType: 0, Level: 60, ClassJob: 15,

			ChocoboInfo: datatypes.ChocoboInfo{
				Head: 1, Body: 1, Feet: 1, Color: 0x12,
			},
			StatusLoopVFX: 0x56,

			Statuses: [30]datatypes.StatusEffect{
				{ActorID: 0xE0000000},
				{ActorID: 0xE0000000},
				{ActorID: 0xE0000000},
				{ActorID: 0xE0000000},
				{ActorID: 0xE0000000},
				{ActorID: 0xE0000000},
				{ActorID: 0xE0000000},
				{ActorID: 0xE0000000},
				{ActorID: 0xE0000000},
				{ActorID: 0xE0000000},
			},
			X:     500,
			Y:     600,
			Z:     700,
			Head:  datatypes.Gear{ModelID: 0x1234, Variant: 0x56, Dye: 0x78},
			Body:  datatypes.Gear{ModelID: 0x1234, Variant: 0x56, Dye: 0x78},
			Hand:  datatypes.Gear{ModelID: 0x1234, Variant: 0x56, Dye: 0x78},
			Leg:   datatypes.Gear{ModelID: 0x1234, Variant: 0x56, Dye: 0x78},
			Foot:  datatypes.Gear{ModelID: 0x1234, Variant: 0x56, Dye: 0x78},
			Ear:   datatypes.Gear{ModelID: 0x1234, Variant: 0x56, Dye: 0x78},
			Neck:  datatypes.Gear{ModelID: 0x1234, Variant: 0x56, Dye: 0x78},
			Wrist: datatypes.Gear{ModelID: 0x1234, Variant: 0x56, Dye: 0x78},
			Ring1: datatypes.Gear{ModelID: 0x1234, Variant: 0x56, Dye: 0x78},
			Ring2: datatypes.Gear{ModelID: 0x1234, Variant: 0x56, Dye: 0x78},

			Name: datatypes.StringToEntityName("Striking Dummy"),

			Model: datatypes.ModelInfo{
				Race: 0xE0, Gender: 0xE0, BodyType: 0xE0, Height: 0xE0, Tribe: 0xE0,
				Face: 0xE0, Hairstyle: 0xE0, HairHighlight: 0xE0, SkinTone: 0xE0,
				OddEyeColor: 0xE0, HairColor: 0xE0, HairHighlightColor: 0xE0,
				FacialFeatures: 0xE0, FacialFeaturesColor: 0xE0, Eyebrows: 0xE0,
				EyeColor: 0xE0, EyeShape: 0xE0, Nose: 0xE0, Jaw: 0xE0, Mouth: 0xE0,
				LipColor: 0xE0, TailLength: 0xE0, TailType: 0xE0, BustSize: 0xE0,
				FacePaintType: 0xE0, FacePaintColor: 0xE0,
			},
			FCTag: datatypes.StringToFCTag(":DDD"),
		}
		b.Data = playerSpawnData

		d.ClassJobData = map[byte]datasheet.ClassJob{
			15: {
				Key:          15,
				Name:         "Dummy",
				Abbreviation: "DUM",
			},
		}

		d.WorldData = datasheet.WorldStore{
			123: {Key: 123, Name: "Foo"},
			456: {Key: 123, Name: "Bar"},
		}

		expectedEntityFields = gstruct.Fields{
			"ID":       Equal(subjectID),
			"Index":    Equal(10),
			"Name":     Equal("Striking Dummy"),
			"TargetID": Equal(uint64(0x9ABC)),
			"OwnerID":  Equal(uint64(0x9ABC)),
			"Level":    Equal(60),
			"ClassJob": Equal(&models.ClassJob{
				ID:           15,
				Name:         "Dummy",
				Abbreviation: "DUM",
			}),
			"IsNpc":    BeFalse(),
			"IsEnemy":  BeFalse(),
			"IsPet":    BeFalse(),
			"BNPCInfo": BeNil(),
			"Resources": Equal(&models.Resources{
				Hp:       29000,
				Mp:       11000,
				MaxHp:    30000,
				MaxMp:    12000,
				Tp:       0,
				LastTick: b.Time,
			}),
			"Location": gstruct.PointTo(gstruct.MatchAllFields(gstruct.Fields{
				"X":           BeNumerically("~", 500, 0.001),
				"Y":           BeNumerically("~", 600, 0.001),
				"Z":           BeNumerically("~", 700, 0.001),
				"Orientation": BeNumerically("~", math.Pi, 0.001),
				"LastUpdated": Equal(b.Time),
			})),
			"LastAction":   BeNil(),
			"Statuses":     BeEmpty(),
			"LockonMarker": Equal(0),
			"CastingInfo":  BeNil(),
		}
	})

	JustBeforeEach(func() {
		rawSpawnData, err := json.Marshal(playerSpawnData)
		Expect(err).ToNot(HaveOccurred())
		expectedEntityFields["RawSpawnJSONData"] = Equal(string(rawSpawnData))
	})

	expectOneEntityToSpawn := func(checkStreamEvents func([]models.StreamEvent)) {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())

		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())

		if checkStreamEvents != nil {
			checkStreamEvents(streamEvents)
		} else {
			Expect(streamEvents).To(BeEmpty())
		}

		Expect(entityEvents).To(HaveLen(1))
		Expect(entityEvents[0].StreamID).To(Equal(streamID))
		Expect(entityEvents[0].EntityID).To(Equal(subjectID))
		eventType, assignable := entityEvents[0].Type.(models.AddEntity)
		Expect(assignable).To(BeTrue())

		Expect(eventType.Entity).To(gstruct.PointTo(gstruct.MatchAllFields(expectedEntityFields)))

		Expect(streams.Map[streamID].EntitiesMap).To(HaveKey(subjectID))
		Expect(streams.Map[streamID].EntitiesMap[subjectID]).To(
			gstruct.PointTo(gstruct.MatchAllFields(expectedEntityFields)),
		)

		Expect(validate.Validate(entityEvents)).To(Succeed())
		Expect(validate.Validate(streams)).To(Succeed())
	}

	It("generates an update to spawn a Player entitty", func() {
		expectOneEntityToSpawn(nil)
	})

	Context("when the spawn is for the current player character", func() {
		BeforeEach(func() {
			b.CurrentID = b.SubjectID
			playerSpawnData.Index = 0
			delete(streams.Map[streamID].EntitiesMap, subjectID)
			expectedEntityFields["Index"] = Equal(0)
		})

		It("generates both a stream event for the world IDs and the spawn entity update", func() {
			expectOneEntityToSpawn(func(streamEvents []models.StreamEvent) {
				Expect(streamEvents).To(HaveLen(1))
				Expect(streamEvents[0].StreamID).To(Equal(streamID))

				eventType, assignable := streamEvents[0].Type.(models.UpdateIDs)
				Expect(assignable).To(BeTrue())

				Expect(eventType.ServerID).To(Equal(2000))
				Expect(eventType.InstanceNum).To(Equal(1000))

				Expect(eventType.CharacterID).To(Equal(subjectID))
				Expect(eventType.CurrentWorld).To(Equal(&models.World{ID: 123, Name: "Foo"}))
				Expect(eventType.HomeWorld).To(Equal(&models.World{ID: 456, Name: "Bar"}))

				Expect(validate.Validate(streamEvents)).To(Succeed())
			})
		})
	})

	Context("with status effects", func() {
		BeforeEach(func() {
			playerSpawnData.Statuses = [30]datatypes.StatusEffect{
				{ID: 0, ActorID: 0xE0000000},
				{ID: 91, Param: 1, Duration: 1, ActorID: 0xE0000000},
				{ID: 92, Param: 2, Duration: 0.5, ActorID: 0xE0000000},
				{ID: 93, Param: 3, Duration: 0.25, ActorID: 0xE0000000},
			}
			d.StatusData = map[uint32]datasheet.Status{
				91: {Key: 91, Name: "Test1", Description: "First"},
				92: {Key: 92, Name: "Test2", Description: "Second"},
				93: {Key: 93, Name: "Test3", Description: "Third"},
			}

			expectedEntityFields["Statuses"] = Equal([]*models.Status{
				nil,
				{
					ID: 91, Param: 1, Name: "Test1", Description: "First",
					StartedTime: b.Time, Duration: time.Unix(1, 0),
					ActorID: 0xE0000000, LastTick: b.Time,
				},
				{
					ID: 92, Param: 2, Name: "Test2", Description: "Second",
					StartedTime: b.Time, Duration: time.Unix(0, 500000000),
					ActorID: 0xE0000000, LastTick: b.Time,
				},
				{
					ID: 93, Param: 3, Name: "Test3", Description: "Third",
					StartedTime: b.Time, Duration: time.Unix(0, 250000000),
					ActorID: 0xE0000000, LastTick: b.Time,
				},
			})
		})

		It("generates an update to spawn the entity with status effects", func() {
			expectOneEntityToSpawn(nil)
		})
	})

	Context("when an entity at the index already exists", func() {
		BeforeEach(func() {
			playerSpawnData.Index = 123
			expectedEntityFields["Index"] = Equal(123)
		})

		It("first generates an update to despawn the old entity and then generates an update to spawn the new entity", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())

			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

			Expect(entityEvents).To(HaveLen(2))

			Expect(entityEvents[0].StreamID).To(Equal(streamID))
			Expect(entityEvents[0].EntityID).To(Equal(removableID))
			removeEvent, assignable := entityEvents[0].Type.(models.RemoveEntity)
			Expect(assignable).To(BeTrue())
			Expect(removeEvent.ID).To(Equal(removableID))

			Expect(entityEvents[1].StreamID).To(Equal(streamID))
			Expect(entityEvents[1].EntityID).To(Equal(subjectID))
			addEvent, assignable := entityEvents[1].Type.(models.AddEntity)
			Expect(assignable).To(BeTrue())
			Expect(addEvent.Entity).To(gstruct.PointTo(gstruct.MatchAllFields(expectedEntityFields)))

			Expect(streams.Map[streamID].EntitiesMap).To(HaveKey(subjectID))
			Expect(streams.Map[streamID].EntitiesMap[subjectID]).To(
				gstruct.PointTo(gstruct.MatchAllFields(expectedEntityFields)),
			)

			Expect(validate.Validate(entityEvents)).To(Succeed())
			Expect(validate.Validate(streams)).To(Succeed())
		})
	})

	Context("when the entity name has decoding errors", func() {
		BeforeEach(func() {
			playerSpawnData.Name = datatypes.StringToEntityName("木\xc5人")
			expectedEntityFields["Name"] = Equal("木人")
		})

		It("generates an update to spawn the entity with the name stripped of invalid characters", func() {
			expectOneEntityToSpawn(nil)
		})
	})

	Context("when the entity is an NPC spawn", func() {
		var npcSpawn *datatypes.NPCSpawn

		BeforeEach(func() {
			d.BNPCData.BNPCNames = testassets.ExpectedBNPCNames
			d.BNPCData.BNPCBases = testassets.ExpectedBNPCBases
			d.BNPCData.ModelCharas = testassets.ExpectedModelCharas
			d.BNPCData.ModelSkeletons = testassets.ExpectedModelSkeletons
			npcSpawn = &datatypes.NPCSpawn{
				PlayerSpawn: *playerSpawnData,
			}
			b.Data = npcSpawn

			npcSpawn.Type = 2
			npcSpawn.Subtype = 4
			npcSpawn.EnemyType = 4

			npcSpawn.BNPCBase = 3
			npcSpawn.BNPCName = 2
			npcSpawn.ModelChara = 878

			expectedEntityFields["IsNpc"] = BeTrue()
			expectedEntityFields["IsEnemy"] = BeTrue()
			expectedEntityFields["IsPet"] = BeFalse()
			npcName := "Ruins Runner"
			npcSize := float64(float32(1.2) * float32(0.2))
			expectedEntityFields["BNPCInfo"] = Equal(&models.NPCInfo{
				NameID:  2,
				BaseID:  3,
				ModelID: 878,
				Name:    &npcName,
				Size:    &npcSize,
			})
			// Since this is an NPC ClassJob shouldn't really be considered
			expectedEntityFields["ClassJob"] = Equal(&models.ClassJob{
				ID: 15,
			})
		})

		JustBeforeEach(func() {
			rawSpawnData, err := json.Marshal(npcSpawn.PlayerSpawn)
			Expect(err).ToNot(HaveOccurred())
			expectedEntityFields["RawSpawnJSONData"] = Equal(string(rawSpawnData))
		})

		It("generates an update to spawn the NPC", func() {
			expectOneEntityToSpawn(nil)
		})

		Context("when the entity is a pet", func() {
			BeforeEach(func() {
				npcSpawn.EnemyType = 0
				npcSpawn.Subtype = 2

				expectedEntityFields["IsNpc"] = BeTrue()
				expectedEntityFields["IsEnemy"] = BeFalse()
				expectedEntityFields["IsPet"] = BeTrue()
			})

			It("generates an update to spawn the friendly entity", func() {
				expectOneEntityToSpawn(nil)
			})
		})

		Context("when the entity is some special entity", func() {
			BeforeEach(func() {
				b.Data = &datatypes.NPCSpawn2{
					PlayerSpawn: npcSpawn.PlayerSpawn,
				}
			})

			It("generates an update to spawn the NPC", func() {
				expectOneEntityToSpawn(nil)
			})
		})
	})

	streamValidationTests(testEnv, false)
})
