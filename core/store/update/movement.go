package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.Movement), newMovementUpdate)
}

func newMovementUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.Movement)

	return locationUpdate{
		streamID:  streamID,
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
	streamID  int
	subjectID uint64

	location models.Location
}

func (u locationUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	return validateEntityUpdate(streams, u.streamID, u.subjectID, u.modifyFunc)
}

func (u locationUpdate) modifyFunc(stream *models.Stream, entity *models.Entity) ([]models.StreamEvent, []models.EntityEvent, error) {
	entity.Location = &u.location

	return nil, []models.EntityEvent{{
		StreamID: u.streamID,
		EntityID: u.subjectID,
		Type: models.UpdateLocation{
			Location: &u.location,
		},
	}}, nil
}
