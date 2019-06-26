package update

import (
	"sort"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.HateRanking), newHateRankingUpdate)
}

// TODO: Add testing
func newHateRankingUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.HateRanking)
	var l []models.HateRanking
	for _, e := range data.Entries[:data.Count] {
		l = append(l, models.HateRanking{
			ActorID: uint64(e.ActorID),
			Hate:    int(e.Hate),
		})
	}
	sort.SliceStable(l, func(i, j int) bool {
		return l[i].Hate > l[j].Hate
	})
	return hateRankingUpdate{
		streamID: streamID,

		hateRankings: l,
	}
}

type hateRankingUpdate struct {
	streamID int

	hateRankings []models.HateRanking
}

func (u hateRankingUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	stream, found := streams.Map[u.streamID]
	if !found {
		return nil, nil, ErrorStreamNotFound
	}

	stream.Enmity.TargetHateRanking = u.hateRankings

	return []models.StreamEvent{models.StreamEvent{
		StreamID: u.streamID,
		Type: models.UpdateEnmity{
			Enmity: stream.Enmity,
		},
	}}, nil, nil
}
