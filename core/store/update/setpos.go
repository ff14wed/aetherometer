package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.SetPos), newSetPosUpdate)
}

func newSetPosUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.SetPos)

	return locationUpdate{
		streamID:  streamID,
		subjectID: uint64(b.SubjectID),
		location: models.Location{
			Orientation: getCanonicalOrientation(uint32(data.Direction), 0x100),
			X:           float64(data.X),
			Y:           float64(data.Y),
			Z:           float64(data.Z),
			LastUpdated: b.Time,
		},
	}
}
