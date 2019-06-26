package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.HateList), newHateListUpdate)
}

// TODO: Add testing
func newHateListUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.HateList)
	var l []models.HateEntry
	for _, e := range data.Entries[:data.Count] {
		l = append(l, models.HateEntry{
			EnemyID:     uint64(e.EnemyID),
			HatePercent: int(e.HatePct),
		})
	}

	return hateListUpdate{
		streamID: streamID,

		hateList: l,
	}
}

type hateListUpdate struct {
	streamID int

	hateList []models.HateEntry
}

func (u hateListUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	stream, found := streams.Map[u.streamID]
	if !found {
		return nil, nil, ErrorStreamNotFound
	}

	stream.Enmity.NearbyEnemyHate = u.hateList

	return []models.StreamEvent{models.StreamEvent{
		StreamID: u.streamID,
		Type: models.UpdateEnmity{
			Enmity: stream.Enmity,
		},
	}}, nil, nil
}
