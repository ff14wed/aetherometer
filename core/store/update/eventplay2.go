package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.EventPlay2), newEventPlay2Update)
}

// TODO: Add testing
func newEventPlay2Update(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.EventPlay2)

	if data.EventID == 0xA0001 {
		switch data.P3 {
		case 2:
			recipeInfo := d.RecipeData.GetInfo(data.P4)
			return craftingInfoUpdate{
				streamID: streamID,
				craftingInfo: &models.CraftingInfo{
					Recipe:            recipeInfo,
					StepNum:           1,
					CurrentCondition:  1,
					PreviousCondition: 1,
				},
			}
		case 4, 6:
			return craftingInfoUpdate{
				streamID:     streamID,
				craftingInfo: nil,
			}
		}
	}
	return nil
}
