package update

import (
	"fmt"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	"go.uber.org/zap"
)

func init() {
	registerIngressHandler(new(datatypes.InitZone), newInitZoneUpdate)
}

func newInitZoneUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.InitZone)

	zap.L().Debug("InitZone", zap.String("data", fmt.Sprintf("%#v", data)))

	var place models.Place
	mapInfos := d.MapData.GetMaps(data.TerritoryTypeID)
	// TODO: Figure out a way to dynamically figure out the current map ID based
	// on character location. Probably do this for movement update handlers.
	if len(mapInfos) > 0 {
		place.MapID = mapInfos[0].Key
	}
	place.TerritoryID = int(data.TerritoryTypeID)
	place.Maps = mapInfos

	// The assumption here is that the current server ID only changes on zone
	// change, but this isn't confirmed. Needs testing especially in cases of
	// high player count in a single zone.
	return placeUpdate{
		streamID:  streamID,
		serverID:  int(b.ServerID),
		currentID: uint64(b.CurrentID),

		place: place,
	}
}

type placeUpdate struct {
	streamID  int
	serverID  int
	currentID uint64

	place models.Place
}

func (u placeUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	stream, found := streams.Map[u.streamID]
	if !found {
		return nil, nil, ErrorStreamNotFound
	}
	if u.place.TerritoryID == stream.Place.TerritoryID {
		return nil, nil, nil
	}

	var streamEvents []models.StreamEvent

	stream.ServerID = u.serverID
	stream.CharacterID = u.currentID
	streamEvents = append(streamEvents, models.StreamEvent{
		StreamID: u.streamID,
		Type: models.UpdateIDs{
			ServerID:    u.serverID,
			CharacterID: u.currentID,
		},
	})

	stream.Place = u.place
	streamEvents = append(streamEvents, models.StreamEvent{
		StreamID: u.streamID,
		Type: models.UpdateMap{
			Place: u.place,
		},
	})

	stream.EntitiesMap = make(map[uint64]*models.Entity)
	entityEvents := []models.EntityEvent{
		models.EntityEvent{
			StreamID: u.streamID,
			Type: models.SetEntities{
				Entities: nil,
			},
		},
	}

	return streamEvents, entityEvents, nil
}
