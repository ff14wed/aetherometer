package update

import (
	"encoding/json"
	"time"
	"unicode/utf8"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.PlayerSpawn), newPlayerSpawnUpdate)
	registerIngressHandler(new(datatypes.NPCSpawn), newNPCSpawnUpdate)
	registerIngressHandler(new(datatypes.NPCSpawn2), newNPCSpawn2Update)
}

func newNPCSpawnUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.NPCSpawn)
	return generateSpawnUpdate(streamID, uint64(b.SubjectID), uint64(b.CurrentID), b.Time, &data.PlayerSpawn, d, true)
}

func newNPCSpawn2Update(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.NPCSpawn2)
	return generateSpawnUpdate(streamID, uint64(b.SubjectID), uint64(b.CurrentID), b.Time, &data.PlayerSpawn, d, true)
}

func newPlayerSpawnUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.PlayerSpawn)
	return generateSpawnUpdate(streamID, uint64(b.SubjectID), uint64(b.CurrentID), b.Time, data, d, false)
}

func generateSpawnUpdate(
	streamID int,
	subjectID uint64,
	currentID uint64,
	now time.Time,
	data *datatypes.PlayerSpawn,
	d *datasheet.Collection,
	isNPCSpawn bool,
) store.Update {
	spawnName := data.Name.String()
	filteredName := make([]rune, 0, len(spawnName))
	for i, r := range spawnName {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(spawnName[i:])
			if size == 1 {
				continue
			}
		}
		filteredName = append(filteredName, r)
	}

	spawnJSONBytes, _ := json.Marshal(data)

	newEntity := models.Entity{
		ID:       uint64(subjectID),
		Name:     string(filteredName),
		Index:    int(data.Index),
		TargetID: data.TargetID,
		OwnerID:  uint64(data.OwnerID),
		Level:    int(data.Level),
		ClassJob: &models.ClassJob{
			ID: int(data.ClassJob),
		},
		IsNpc:   isNPCSpawn,
		IsEnemy: (data.EnemyType != 0),
		IsPet:   (data.EnemyType == 0 && data.Subtype == 2),

		Resources: &models.Resources{
			Hp:       int(data.CurrentHP),
			Mp:       int(data.CurrentMP),
			MaxHp:    int(data.MaxHP),
			MaxMp:    int(data.MaxMP),
			LastTick: now,
		},
		Location: &models.Location{
			Orientation: getCanonicalOrientation(uint32(data.Direction), 0x10000),
			X:           float64(data.X),
			Z:           float64(data.Z),
			Y:           float64(data.Y),
			LastUpdated: now,
		},
		RawSpawnJSONData: string(spawnJSONBytes),
	}

	if isNPCSpawn {
		newEntity.BNPCInfo = &models.NPCInfo{
			NameID:  int(data.BNPCName),
			BaseID:  int(data.BNPCBase),
			ModelID: int(data.ModelChara),
		}
		if bNPCInfo := d.BNPCData.GetBNPCInfo(
			data.BNPCName,
			data.BNPCBase,
			uint32(data.ModelChara),
		); bNPCInfo != nil {
			size := float64(bNPCInfo.Size)
			newEntity.BNPCInfo.Name = &bNPCInfo.Name
			newEntity.BNPCInfo.Size = &size
			newEntity.BNPCInfo.Error = int(bNPCInfo.Error)
		}
	} else if data.ClassJob > 0 {
		if cj, found := d.ClassJobData[data.ClassJob]; found {
			newEntity.ClassJob.Name = cj.Name
			newEntity.ClassJob.Abbreviation = cj.Abbreviation
		}
	}

	var statusListLength byte

	for i := byte(0); i < 30; i++ {
		if data.Statuses[i].ID != 0 {
			statusListLength = i + 1
		}
	}

	newEntity.Statuses = make([]*models.Status, statusListLength)

	for i := byte(0); i < statusListLength; i++ {
		status := data.Statuses[i]
		if status.ID == 0 {
			continue
		}
		var name, description string

		if statusData, found := d.StatusData[uint32(status.ID)]; found {
			name = statusData.Name
			description = statusData.Description
		}
		newEntity.Statuses[i] = &models.Status{
			ID:          int(status.ID),
			Param:       int(status.Param),
			Name:        name,
			Description: description,
			StartedTime: now,
			Duration:    getTimeForDuration(status.Duration),
			ActorID:     uint64(status.ActorID),
			LastTick:    now,
		}
	}

	var (
		homeWorld, currentWorld models.World
		isWorldSet              bool
	)

	if subjectID == currentID {
		homeWorld = d.WorldData.Lookup(int(data.HomeWorld))
		currentWorld = d.WorldData.Lookup(int(data.CurrentWorld))
		isWorldSet = true
	}

	return entitySpawnUpdate{
		streamID:  streamID,
		subjectID: subjectID,

		isWorldSet:   isWorldSet,
		homeWorld:    homeWorld,
		currentWorld: currentWorld,

		entity: newEntity,
	}
}

type entitySpawnUpdate struct {
	streamID  int
	subjectID uint64

	isWorldSet   bool
	homeWorld    models.World
	currentWorld models.World

	entity models.Entity
}

func (u entitySpawnUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	stream, found := streams.Map[u.streamID]
	if !found {
		return nil, nil, ErrorStreamNotFound
	}

	var (
		streamEvents []models.StreamEvent
		entityEvents []models.EntityEvent
	)

	for key, ent := range stream.EntitiesMap {
		if ent == nil || ent.Index != u.entity.Index {
			continue
		}
		stream.EntitiesMap[key] = nil
		entityEvents = append(entityEvents, models.EntityEvent{
			StreamID: u.streamID,
			EntityID: key,
			Type: models.RemoveEntity{
				ID: key,
			},
		})
		break
	}

	entityClone := u.entity.Clone()
	entityEvents = append(entityEvents, models.EntityEvent{
		StreamID: u.streamID,
		EntityID: u.subjectID,
		Type: models.AddEntity{
			Entity: &entityClone,
		},
	})

	stream.EntitiesMap[u.subjectID] = &u.entity

	if u.isWorldSet {
		stream.HomeWorld = u.homeWorld
		stream.CurrentWorld = u.currentWorld

		streamEvents = append(streamEvents, models.StreamEvent{
			StreamID: u.streamID,
			Type: models.UpdateIDs{
				ServerID:    stream.ServerID,
				InstanceNum: stream.InstanceNum,

				CharacterID:  stream.CharacterID,
				HomeWorld:    &u.homeWorld,
				CurrentWorld: &u.currentWorld,
			},
		})
	}

	return streamEvents, entityEvents, nil
}
