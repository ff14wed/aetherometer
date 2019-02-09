package update

import (
	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/xivnet/v2"
	"github.com/ff14wed/xivnet/v2/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.AoEAction8), newAoEAction8Update)
	registerIngressHandler(new(datatypes.AoEAction16), newAoEAction16Update)
	registerIngressHandler(new(datatypes.AoEAction24), newAoEAction24Update)
	registerIngressHandler(new(datatypes.AoEAction32), newAoEAction32Update)
}

func newAoEAction8Update(pid int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.AoEAction8)

	action := actionFromHeader(data.ActionHeader, d, b.Header.Time)

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
		pid:       pid,
		subjectID: uint64(b.Header.SubjectID),
		action:    action,
	}
}

func newAoEAction16Update(pid int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.AoEAction16)

	action := actionFromHeader(data.ActionHeader, d, b.Header.Time)

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
		pid:       pid,
		subjectID: uint64(b.Header.SubjectID),
		action:    action,
	}
}

func newAoEAction24Update(pid int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.AoEAction24)

	action := actionFromHeader(data.ActionHeader, d, b.Header.Time)

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
		pid:       pid,
		subjectID: uint64(b.Header.SubjectID),
		action:    action,
	}
}

func newAoEAction32Update(pid int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.AoEAction32)

	action := actionFromHeader(data.ActionHeader, d, b.Header.Time)

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
		pid:       pid,
		subjectID: uint64(b.Header.SubjectID),
		action:    action,
	}
}
