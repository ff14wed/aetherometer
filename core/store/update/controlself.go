package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.ControlSelf), newControlSelfUpdate)
}

func newControlSelfUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.ControlSelf)

	switch data.Type {
	case 0x101:
		return removeEntityUpdate{
			streamID:  streamID,
			subjectID: uint64(data.P3),
		}
	}
	return nil
}
