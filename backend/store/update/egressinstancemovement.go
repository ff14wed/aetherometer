package update

import (
	"math"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerEgressHandler(new(datatypes.EgressInstanceMovement), newEgressInstanceMovementUpdate)
}

func newEgressInstanceMovementUpdate(pid int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.EgressInstanceMovement)

	return locationUpdate{
		pid:       pid,
		subjectID: uint64(b.SubjectID),
		location: models.Location{
			Orientation: math.Pi + float64(data.Direction),
			X:           float64(data.X),
			Y:           float64(data.Y),
			Z:           float64(data.Z),
			LastUpdated: b.Time,
		},
	}
}
