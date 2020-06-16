package update

import (
	"time"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.UpdateStatuses), newUpdateStatusesUpdate)
	registerIngressHandler(new(datatypes.UpdateStatusesEureka), newUpdateStatusesEurekaUpdate)
}

// TODO: Add testing
func processUpdatedStatuses(
	statuses [30]datatypes.StatusEffect,
	d *datasheet.Collection,
	time time.Time,
) (map[int]models.Status, byte) {
	updatedStatuses := make(map[int]models.Status)

	var statusListLength byte

	for i := byte(0); i < 30; i++ {
		e := statuses[i]
		if e.ID == 0 {
			continue
		}

		var name, description string

		if statusData, found := d.StatusData[uint32(e.ID)]; found {
			name = statusData.Name
			description = statusData.Description
		}
		updatedStatuses[int(i)] = models.Status{
			ID:          int(e.ID),
			Param:       int(e.Param),
			Name:        name,
			Description: description,
			StartedTime: time,
			Duration:    getTimeForDuration(e.Duration),
			ActorID:     uint64(e.ActorID),
			LastTick:    time,
		}
		if i >= statusListLength {
			statusListLength = i + 1
		}
	}
	return updatedStatuses, statusListLength
}

func newUpdateStatusesUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.UpdateStatuses)

	updatedStatuses, statusListLength := processUpdatedStatuses(data.Statuses, d, b.Time)
	return updateStatusesUpdate{
		streamID:  streamID,
		subjectID: uint64(b.SubjectID),

		statusListLength: byte(statusListLength),
		statuses:         updatedStatuses,
		resources: models.Resources{
			Hp:       int(data.CurrentHP),
			Mp:       int(data.CurrentMP),
			Tp:       int(data.CurrentTP),
			MaxHp:    int(data.MaxHP),
			MaxMp:    int(data.MaxMP),
			LastTick: b.Time,
		},
	}
}

func newUpdateStatusesEurekaUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.UpdateStatusesEureka)

	updatedStatuses, statusListLength := processUpdatedStatuses(data.Statuses, d, b.Time)
	return updateStatusesUpdate{
		streamID:  streamID,
		subjectID: uint64(b.SubjectID),

		statusListLength: byte(statusListLength),
		statuses:         updatedStatuses,
		resources: models.Resources{
			Hp:       int(data.CurrentHP),
			Mp:       int(data.CurrentMP),
			Tp:       int(data.CurrentTP),
			MaxHp:    int(data.MaxHP),
			MaxMp:    int(data.MaxMP),
			LastTick: b.Time,
		},
	}
}

type updateStatusesUpdate struct {
	streamID  int
	subjectID uint64

	statusListLength byte
	statuses         map[int]models.Status
	resources        models.Resources
}

func (u updateStatusesUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	return validateEntityUpdate(streams, u.streamID, u.subjectID, u.modifyFunc)
}

func (u updateStatusesUpdate) modifyFunc(stream *models.Stream, entity *models.Entity) ([]models.StreamEvent, []models.EntityEvent, error) {
	resourcesClone := u.resources
	entity.Resources = &u.resources

	if len(entity.Statuses) < int(u.statusListLength) {
		diff := int(u.statusListLength) - len(entity.Statuses)
		entity.Statuses = append(entity.Statuses, make([]*models.Status, diff)...)
	}

	var statusEvents []models.EntityEvent
	for i, s := range entity.Statuses {
		updatedStatus := u.statuses[i]
		if updatedStatus.ID == 0 && s != nil {
			entity.Statuses[i] = nil
			statusEvents = append(statusEvents, models.EntityEvent{
				StreamID: u.streamID,
				EntityID: u.subjectID,
				Type:     models.RemoveStatus{Index: i},
			})
		} else if updatedStatus.ID > 0 {
			if entity.Statuses[i] != nil {
				updatedStatus.StartedTime = entity.Statuses[i].StartedTime
			}
			entity.Statuses[i] = &updatedStatus
			statusEvents = append(statusEvents, models.EntityEvent{
				StreamID: u.streamID,
				EntityID: u.subjectID,
				Type:     models.UpsertStatus{Index: i, Status: &updatedStatus},
			})
		}
	}
	entityEvents := append(statusEvents, models.EntityEvent{
		StreamID: u.streamID,
		EntityID: u.subjectID,
		Type:     models.UpdateResources{Resources: &resourcesClone},
	})

	return nil, entityEvents, nil
}
