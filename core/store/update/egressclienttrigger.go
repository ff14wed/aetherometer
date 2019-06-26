package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerEgressHandler(new(datatypes.EgressClientTrigger), newEgressClientTriggerUpdate)
}

// TODO: Add testing
func newEgressClientTriggerUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.EgressClientTrigger)

	switch data.Type {
	case 0x3:
		return targetUpdate{
			streamID:  streamID,
			subjectID: uint64(b.SubjectID),

			targetID: uint64(data.P1),
		}
	}
	return nil
}
