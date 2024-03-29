package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.ControlTarget), newControlTargetUpdate)
}

func newControlTargetUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.ControlTarget)

	switch data.Type {
	case 0x32:
		return targetUpdate{
			streamID:  streamID,
			subjectID: uint64(b.SubjectID),

			targetID: uint64(data.TargetID),
		}
	}
	return nil
}

type targetUpdate struct {
	streamID  int
	subjectID uint64

	targetID uint64
}

func (u targetUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	return validateEntityUpdate(streams, u.streamID, u.subjectID, u.modifyFunc)
}

func (u targetUpdate) modifyFunc(stream *models.Stream, entity *models.Entity) ([]models.StreamEvent, []models.EntityEvent, error) {
	entity.TargetID = u.targetID

	return nil, []models.EntityEvent{{
		StreamID: u.streamID,
		EntityID: u.subjectID,
		Type: models.UpdateTarget{
			TargetID: u.targetID,
		},
	}}, nil
}
