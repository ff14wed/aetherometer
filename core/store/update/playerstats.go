package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.PlayerStats), newPlayerStatsUpdate)
}

func newPlayerStatsUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.PlayerStats)

	return playerStatsUpdate{
		streamID: streamID,

		stats: models.Stats{
			Strength:            int(data.Strength),
			Dexterity:           int(data.Dexterity),
			Vitality:            int(data.Vitality),
			Intelligence:        int(data.Intelligence),
			Mind:                int(data.Mind),
			Piety:               int(data.Piety),
			Hp:                  int(data.HP),
			Mp:                  int(data.MP),
			Tp:                  int(data.TP),
			Gp:                  int(data.GP),
			Cp:                  int(data.CP),
			Delay:               int(data.Delay),
			Tenacity:            int(data.Tenacity),
			AttackPower:         int(data.AttackPower),
			Defense:             int(data.Defense),
			DirectHitRate:       int(data.DirectHitRate),
			Evasion:             int(data.Evasion),
			MagicDefense:        int(data.MagicDefense),
			CriticalHit:         int(data.CriticalHit),
			AttackMagicPotency:  int(data.AttackMagicPotency),
			HealingMagicPotency: int(data.HealingMagicPotency),
			ElementalBonus:      int(data.ElementalBonus),
			Determination:       int(data.Determination),
			SkillSpeed:          int(data.SkillSpeed),
			SpellSpeed:          int(data.SpellSpeed),
			Haste:               int(data.Haste),
			Craftsmanship:       int(data.Craftsmanship),
			Control:             int(data.Control),
			Gathering:           int(data.Gathering),
			Perception:          int(data.Perception),
		},
	}
}

type playerStatsUpdate struct {
	streamID int

	stats models.Stats
}

func (u playerStatsUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	stream, found := streams.Map[u.streamID]
	if !found {
		return nil, nil, ErrorStreamNotFound
	}

	stream.Stats = &u.stats

	return []models.StreamEvent{
		{
			StreamID: u.streamID,
			Type:     models.UpdateStats{Stats: &u.stats},
		},
	}, nil, nil
}
