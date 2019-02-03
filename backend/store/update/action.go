package update

import (
	"fmt"
	"time"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/xivnet"
	"github.com/ff14wed/xivnet/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.Action), newActionUpdate)
}

func actionFromHeader(h datatypes.ActionHeader, d *datasheet.Collection, t time.Time) models.Action {
	id := int(h.ActionID)
	actionName := fmt.Sprintf("Unknown_%x", h.ActionIDName)
	if actionData, found := d.ActionData[(uint32(h.ActionIDName))]; found {
		actionName = actionData.Name
	}
	return models.Action{
		TargetID:          uint64(h.TargetID),
		Name:              actionName,
		GlobalCounter:     int(h.GlobalCounter),
		AnimationLockTime: float64(h.AnimationLockTime),
		HiddenAnimation:   int(h.HiddenAnimation),
		Location: models.Location{
			Orientation: getCanonicalOrientation(uint32(h.Direction), 0xFFFF),
			LastUpdated: t,
		},
		ID:                id,
		Variation:         3,
		EffectDisplayType: 4,
		UseTime:           t,
	}
}

func processActionEffects(effectsList []datatypes.ActionEffects, targets []uint64) []models.ActionEffect {
	if len(effectsList) != len(targets) {
		// This error should never happen due to bad data, only bad code
		panic(fmt.Errorf("effects list length (%d) != target list length (%d)", len(effectsList), len(targets)))
	}
	var actionEffects []models.ActionEffect
	for i, effects := range effectsList {
		target := targets[i]
		for _, e := range effects {
			if e.Type == 0 {
				break
			}
			actionEffects = append(actionEffects, models.ActionEffect{
				TargetID:        target,
				Type:            int(e.Type),
				HitSeverity:     int(e.HitSeverity),
				Param:           int(e.P3),
				BonusPercent:    int(e.Percentage),
				ValueMultiplier: int(e.Multiplier),
				Flags:           int(e.Flags),
				Value:           int(e.Damage),
			})
		}
	}
	return actionEffects
}

func newActionUpdate(pid int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.Action)

	action := actionFromHeader(data.ActionHeader, d, b.Header.Time)

	var actionEffects []models.ActionEffect

	if data.NumAffected > 0 {
		actionEffects = processActionEffects(
			[]datatypes.ActionEffects{data.Effects},
			[]uint64{uint64(data.TargetID2)},
		)
	}
	action.Effects = actionEffects
	action.EffectFlags = int(data.EffectFlags)

	return actionUpdate{
		pid:       pid,
		subjectID: uint64(b.Header.SubjectID),

		action: action,
	}
}

type actionUpdate struct {
	pid       int
	subjectID uint64

	action models.Action
}

func (u actionUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	stream, found := streams.Map[u.pid]
	if !found {
		return nil, nil, ErrorStreamNotFound
	}
	entity, found := stream.EntitiesMap[u.subjectID]
	if !found {
		return nil, nil, ErrorEntityNotFound
	}
	if entity == nil {
		return nil, nil, nil
	}

	entity.LastAction = &u.action

	entityEvents := []models.EntityEvent{models.EntityEvent{
		StreamID: u.pid,
		EntityID: u.subjectID,
		Type: models.UpdateLastAction{
			Action: *entity.LastAction,
		},
	}}

	if entity.CastingInfo != nil {
		entity.CastingInfo = nil
		entityEvents = append(entityEvents, models.EntityEvent{
			StreamID: u.pid,
			EntityID: u.subjectID,
			Type: models.UpdateCastingInfo{
				CastingInfo: nil,
			},
		})
	}

	return nil, entityEvents, nil
}
