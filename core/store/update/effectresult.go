package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.EffectResult), newEffectResultUpdate)
}

// TODO: Add testing
func newEffectResultUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.EffectResult)

	addedStatuses := make(map[int]models.Status)

	var statusListLength byte

	for i := byte(0); i < data.Count; i++ {
		e := data.Entries[i]
		var name, description string

		if statusData, found := d.StatusData[uint32(e.EffectID)]; found {
			name = statusData.Name
			description = statusData.Description
		}
		addedStatuses[int(e.Index)] = models.Status{
			ID:          int(e.EffectID),
			Param:       int(e.Param),
			Name:        name,
			Description: description,
			StartedTime: b.Time,
			Duration:    getTimeForDuration(e.Duration),
			ActorID:     uint64(e.ActorID),
			LastTick:    b.Time,
		}
		if e.Index >= statusListLength {
			statusListLength = e.Index + 1
		}
	}

	return effectResultUpdate{
		streamID:  streamID,
		subjectID: uint64(b.SubjectID),

		statusListLength: statusListLength,
		statuses:         addedStatuses,
		resources: models.Resources{
			Hp:       int(data.CurrentHP),
			Mp:       int(data.CurrentMP),
			Tp:       int(data.CurrentTP),
			MaxHP:    int(data.MaxHP),
			MaxMP:    int(data.MaxMP),
			LastTick: b.Time,
		},
	}
}

type effectResultUpdate struct {
	streamID  int
	subjectID uint64

	statusListLength byte
	statuses         map[int]models.Status
	resources        models.Resources
}

func (u effectResultUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	return validateEntityUpdate(streams, u.streamID, u.subjectID, u.modifyFunc)
}

func (u effectResultUpdate) modifyFunc(stream *models.Stream, entity *models.Entity) ([]models.StreamEvent, []models.EntityEvent, error) {
	entity.Resources = u.resources

	if len(entity.Statuses) <= int(u.statusListLength) {
		diff := int(u.statusListLength) - len(entity.Statuses)
		entity.Statuses = append(entity.Statuses, make([]*models.Status, diff)...)
	}

	var statusEvents []models.EntityEvent
	for i, s := range u.statuses {
		entity.Statuses[i] = &s
		statusEvents = append(statusEvents, models.EntityEvent{
			StreamID: u.streamID,
			EntityID: u.subjectID,
			Type:     models.UpsertStatus{Index: i, Status: s},
		})
	}

	return nil, append([]models.EntityEvent{
		models.EntityEvent{
			StreamID: u.streamID,
			EntityID: u.subjectID,
			Type:     models.UpdateResources{Resources: entity.Resources},
		},
	}, statusEvents...), nil
}
