package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.EventPlay64), newEventPlay64Update)
}

// TODO: Add testing
func newEventPlay64Update(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.EventPlay64)

	craftState, err := datatypes.UnmarshalCraftState(data)
	if err != nil {
		return nil
	}

	flags := craftState.Flags
	completed := (flags & 0x4) > 0
	failed := (flags&0x8) > 0 && !completed

	action := d.ActionData.GetAction(craftState.CraftAction)
	return craftingInfoUpdate{
		streamID: streamID,
		craftingInfo: &models.CraftingInfo{
			LastCraftActionID:   int(craftState.CraftAction),
			LastCraftActionName: action.Name,
			StepNum:             int(craftState.StepNum),
			Progress:            int(craftState.Progress),
			ProgressDelta:       int(craftState.ProgressDelta),
			Quality:             int(craftState.Quality),
			QualityDelta:        int(craftState.QualityDelta),
			HqChance:            int(craftState.HQChance),
			Durability:          int(craftState.Durability),
			DurabilityDelta:     int(craftState.DurabilityDelta),
			CurrentCondition:    int(craftState.CurrentCondition),
			PreviousCondition:   int(craftState.PreviousCondition),
			Completed:           completed,
			Failed:              failed,
		},
	}
}

type craftingInfoUpdate struct {
	streamID int

	craftingInfo *models.CraftingInfo
}

func (u craftingInfoUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	stream, found := streams.Map[u.streamID]
	if !found {
		return nil, nil, ErrorStreamNotFound
	}

	if u.craftingInfo != nil && (u.craftingInfo.Recipe == nil || u.craftingInfo.Recipe.ID == 0) {
		if stream.CraftingInfo != nil {
			u.craftingInfo.Recipe = stream.CraftingInfo.Recipe
		}
	}

	stream.CraftingInfo = u.craftingInfo
	return []models.StreamEvent{
		{
			StreamID: u.streamID,
			Type: models.UpdateCraftingInfo{
				CraftingInfo: u.craftingInfo,
			},
		},
	}, nil, nil
}
