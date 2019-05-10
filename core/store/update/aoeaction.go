package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.AoEAction8), newAoEAction8Update)
	registerIngressHandler(new(datatypes.AoEAction16), newAoEAction16Update)
	registerIngressHandler(new(datatypes.AoEAction24), newAoEAction24Update)
	registerIngressHandler(new(datatypes.AoEAction32), newAoEAction32Update)
}

func newAoEAction8Update(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.AoEAction8)

	action := actionFromHeader(data.ActionHeader, d, b.Time)

	var actionEffects []models.ActionEffect
	numAffected := data.NumAffected
	if numAffected > 0 {
		actionEffects = processActionEffects(
			data.EffectsList[:numAffected],
			data.Targets[:numAffected],
		)
	}
	action.Effects = actionEffects

	return actionUpdate{
		streamID:  streamID,
		subjectID: uint64(b.SubjectID),
		action:    action,
	}
}

func newAoEAction16Update(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.AoEAction16)

	action := actionFromHeader(data.ActionHeader, d, b.Time)

	var actionEffects []models.ActionEffect
	numAffected := data.NumAffected
	if numAffected > 0 {
		actionEffects = processActionEffects(
			data.EffectsList[:numAffected],
			data.Targets[:numAffected],
		)
	}
	action.Effects = actionEffects

	return actionUpdate{
		streamID:  streamID,
		subjectID: uint64(b.SubjectID),
		action:    action,
	}
}

func newAoEAction24Update(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.AoEAction24)

	action := actionFromHeader(data.ActionHeader, d, b.Time)

	var actionEffects []models.ActionEffect
	numAffected := data.NumAffected
	if numAffected > 0 {
		actionEffects = processActionEffects(
			data.EffectsList[:numAffected],
			data.Targets[:numAffected],
		)
	}
	action.Effects = actionEffects

	return actionUpdate{
		streamID:  streamID,
		subjectID: uint64(b.SubjectID),
		action:    action,
	}
}

func newAoEAction32Update(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.AoEAction32)

	action := actionFromHeader(data.ActionHeader, d, b.Time)

	var actionEffects []models.ActionEffect
	numAffected := data.NumAffected
	if numAffected > 0 {
		actionEffects = processActionEffects(
			data.EffectsList[:numAffected],
			data.Targets[:numAffected],
		)
	}
	action.Effects = actionEffects

	return actionUpdate{
		streamID:  streamID,
		subjectID: uint64(b.SubjectID),
		action:    action,
	}
}
