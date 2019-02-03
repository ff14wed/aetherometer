package update

import (
	"fmt"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/xivnet"
	"github.com/ff14wed/xivnet/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.Casting), newCastingUpdate)
}

func newCastingUpdate(pid int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.Casting)

	info := &models.CastingInfo{
		ActionID:  int(data.ActionID),
		StartTime: b.Header.Time,
		CastTime:  getTimeForDuration(data.CastTime),
		TargetID:  uint64(data.TargetID),
		Location: models.Location{
			Orientation: float64(data.Direction),
			X:           float64(data.Position.X.Float()),
			Y:           float64(data.Position.Y.Float()),
			Z:           float64(data.Position.Z.Float()),
		},
	}

	info.ActionName = fmt.Sprintf("Unknown_%x", data.ActionIDName)
	if actionInfo, ok := d.ActionData[uint32(data.ActionIDName)]; ok {
		info.ActionName = actionInfo.Name
	}

	if actionInfo, ok := d.ActionData[uint32(data.ActionID)]; ok {
		info.CastType = int(actionInfo.CastType)
		info.EffectRange = int(actionInfo.EffectRange)
		info.XAxisModifier = int(actionInfo.XAxisModifier)
		info.Omen = actionInfo.Omen
	}

	return castingUpdate{
		pid:       pid,
		subjectID: uint64(b.Header.SubjectID),

		castingInfo: info,
	}
}

type castingUpdate struct {
	pid       int
	subjectID uint64

	castingInfo *models.CastingInfo
}

func (u castingUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
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

	entity.CastingInfo = u.castingInfo

	return nil, []models.EntityEvent{models.EntityEvent{
		StreamID: u.pid,
		EntityID: u.subjectID,
		Type: models.UpdateCastingInfo{
			CastingInfo: u.castingInfo,
		},
	}}, nil
}
