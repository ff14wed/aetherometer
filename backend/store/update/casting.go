package update

import (
	"fmt"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.Casting), newCastingUpdate)
}

func newCastingUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.Casting)

	info := &models.CastingInfo{
		ActionID:  int(data.ActionID),
		StartTime: b.Time,
		CastTime:  getTimeForDuration(data.CastTime),
		TargetID:  uint64(data.TargetID),
		Location: models.Location{
			Orientation: float64(data.Direction),
			X:           float64(data.Position.X.Float()),
			Y:           float64(data.Position.Y.Float()),
			Z:           float64(data.Position.Z.Float()),
			LastUpdated: b.Time,
		},
	}

	info.ActionName = fmt.Sprintf("Unknown_%x", data.ActionIDName)

	actionInfo := d.ActionData.GetAction(uint32(data.ActionIDName))
	if actionInfo.Key != 0 {
		info.ActionName = actionInfo.Name
	}

	actionInfo = d.ActionData.GetAction(data.ActionID)
	if actionInfo.Key != 0 {
		info.CastType = int(actionInfo.CastType)
		info.EffectRange = int(actionInfo.EffectRange)
		info.XAxisModifier = int(actionInfo.XAxisModifier)
		info.Omen = actionInfo.Omen
	}

	return castingUpdate{
		streamID:  streamID,
		subjectID: uint64(b.SubjectID),

		castingInfo: info,
	}
}

type castingUpdate struct {
	streamID  int
	subjectID uint64

	castingInfo *models.CastingInfo
}

func (u castingUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	return validateEntityUpdate(streams, u.streamID, u.subjectID, u.modifyFunc)
}

func (u castingUpdate) modifyFunc(stream *models.Stream, entity *models.Entity) ([]models.StreamEvent, []models.EntityEvent, error) {
	entity.CastingInfo = u.castingInfo

	return nil, []models.EntityEvent{models.EntityEvent{
		StreamID: u.streamID,
		EntityID: u.subjectID,
		Type: models.UpdateCastingInfo{
			CastingInfo: u.castingInfo,
		},
	}}, nil
}
