package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.EquipChange), newEquipChangeUpdate)
}

func newEquipChangeUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.EquipChange)

	var className, classAbbrev string

	if cj, found := d.ClassJobData[data.ClassJob]; found {
		className = cj.Name
		classAbbrev = cj.Abbreviation
	}

	return equipChangeUpdate{
		streamID:  streamID,
		subjectID: uint64(b.SubjectID),

		level: int(data.Level),

		classJob: models.ClassJob{
			ID:           int(data.ClassJob),
			Name:         className,
			Abbreviation: classAbbrev,
		},
	}
}

type equipChangeUpdate struct {
	streamID  int
	subjectID uint64

	level    int
	classJob models.ClassJob
}

func (u equipChangeUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	return validateEntityUpdate(streams, u.streamID, u.subjectID, u.modifyFunc)
}

func (u equipChangeUpdate) modifyFunc(stream *models.Stream, entity *models.Entity) ([]models.StreamEvent, []models.EntityEvent, error) {
	entity.ClassJob = u.classJob
	entity.Level = u.level

	return nil, []models.EntityEvent{models.EntityEvent{
		StreamID: u.streamID,
		EntityID: u.subjectID,
		Type: models.UpdateClass{
			ClassJob: u.classJob,
			Level:    u.level,
		},
	}}, nil
}
