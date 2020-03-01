package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.UpdateHPMPTP), newHPMPTPUpdate)
}

func newHPMPTPUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.UpdateHPMPTP)

	return hpmptpUpdate{
		streamID:  streamID,
		subjectID: uint64(b.SubjectID),

		resources: models.Resources{
			Hp:       int(data.HP),
			Mp:       int(data.MP),
			LastTick: b.Time,
		},
	}
}

type hpmptpUpdate struct {
	streamID  int
	subjectID uint64

	resources models.Resources
}

func (u hpmptpUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	return validateEntityUpdate(streams, u.streamID, u.subjectID, u.modifyFunc)
}

func (u hpmptpUpdate) modifyFunc(stream *models.Stream, entity *models.Entity) ([]models.StreamEvent, []models.EntityEvent, error) {
	entity.Resources.Hp = u.resources.Hp
	if entity.ClassJob.ID < 8 || entity.ClassJob.ID > 18 {
		entity.Resources.Mp = u.resources.Mp
	}
	entity.Resources.LastTick = u.resources.LastTick

	return nil, []models.EntityEvent{
		models.EntityEvent{
			StreamID: u.streamID,
			EntityID: u.subjectID,
			Type:     models.UpdateResources{Resources: entity.Resources},
		},
	}, nil
}
