package update

import (
	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/xivnet/v2"
	"github.com/ff14wed/xivnet/v2/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.EquipChange), newEquipChangeUpdate)
}

func newEquipChangeUpdate(pid int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.EquipChange)

	return equipChangeUpdate{
		pid:       pid,
		subjectID: uint64(b.Header.SubjectID),

		classJob: models.ClassJob{
			ID: int(data.ClassJob),
		},
	}
}

type equipChangeUpdate struct {
	pid       int
	subjectID uint64

	classJob models.ClassJob
}

func (u equipChangeUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	return validateEntityUpdate(streams, u.pid, u.subjectID, u.modifyFunc)
}

func (u equipChangeUpdate) modifyFunc(stream *models.Stream, entity *models.Entity) ([]models.StreamEvent, []models.EntityEvent, error) {
	entity.ClassJob = u.classJob

	return nil, []models.EntityEvent{models.EntityEvent{
		StreamID: u.pid,
		EntityID: u.subjectID,
		Type: models.UpdateClass{
			ClassJob: u.classJob,
		},
	}}, nil
}
