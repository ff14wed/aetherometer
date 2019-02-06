package update

import (
	"math"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/xivnet"
	"github.com/ff14wed/xivnet/datatypes"
)

func init() {
	registerEgressHandler(new(datatypes.MyMovement2), newMovement2Update)
}

func newMovement2Update(pid int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.MyMovement2)

	return locationUpdate{
		pid:       pid,
		subjectID: uint64(b.Header.SubjectID),
		location: models.Location{
			Orientation: math.Pi + float64(data.Direction),
			X:           float64(data.X),
			Y:           float64(data.Y),
			Z:           float64(data.Z),
			LastUpdated: b.Header.Time,
		},
	}
}
