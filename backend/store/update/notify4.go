package update

import (
	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/xivnet/v2"
	"github.com/ff14wed/xivnet/v2/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.Notify4), newNotify4Update)
}

func newNotify4Update(pid int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.Notify4)

	switch data.Type {
	case 0x32:
		return targetUpdate{
			pid:       pid,
			subjectID: uint64(b.Header.SubjectID),

			targetID: uint64(data.TargetID),
		}
	}
	return nil
}

type targetUpdate struct {
	pid       int
	subjectID uint64

	targetID uint64
}

func (u targetUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	return validateEntityUpdate(streams, u.pid, u.subjectID, u.modifyFunc)
}

func (u targetUpdate) modifyFunc(stream *models.Stream, entity *models.Entity) ([]models.StreamEvent, []models.EntityEvent, error) {
	entity.TargetID = u.targetID

	return nil, []models.EntityEvent{models.EntityEvent{
		StreamID: u.pid,
		EntityID: u.subjectID,
		Type: models.UpdateTarget{
			TargetID: u.targetID,
		},
	}}, nil
}
