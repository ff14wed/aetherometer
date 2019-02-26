package update

import (
	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.Movement), newMovementUpdate)
}

func newMovementUpdate(pid int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.Movement)

	return locationUpdate{
		pid:       pid,
		subjectID: uint64(b.SubjectID),
		location: models.Location{
			Orientation: getCanonicalOrientation(uint32(data.Direction), 0x100),
			X:           float64(data.Position.X.Float()),
			Y:           float64(data.Position.Y.Float()),
			Z:           float64(data.Position.Z.Float()),
			LastUpdated: b.Time,
		},
	}
}

type locationUpdate struct {
	pid       int
	subjectID uint64

	location models.Location
}

func (u locationUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	return validateEntityUpdate(streams, u.pid, u.subjectID, u.modifyFunc)
}

func (u locationUpdate) modifyFunc(stream *models.Stream, entity *models.Entity) ([]models.StreamEvent, []models.EntityEvent, error) {
	entity.Location = u.location

	return nil, []models.EntityEvent{models.EntityEvent{
		StreamID: u.pid,
		EntityID: u.subjectID,
		Type: models.UpdateLocation{
			Location: u.location,
		},
	}}, nil
}
