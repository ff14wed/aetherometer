package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.RemoveEntity), newRemoveEntityUpdate)
}

func newRemoveEntityUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.RemoveEntity)

	return removeEntityUpdate{
		streamID:  streamID,
		subjectID: uint64(data.ID),
	}
}

type removeEntityUpdate struct {
	streamID  int
	subjectID uint64
}

func (u removeEntityUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	stream, found := streams.Map[u.streamID]
	if !found {
		return nil, nil, ErrorStreamNotFound
	}
	stream.EntitiesMap[u.subjectID] = nil

	return nil, []models.EntityEvent{{
		StreamID: u.streamID,
		EntityID: u.subjectID,
		Type: models.RemoveEntity{
			ID: u.subjectID,
		},
	}}, nil
}
