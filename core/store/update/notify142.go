package update

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

func init() {
	registerIngressHandler(new(datatypes.Notify142), newNotify142Update)
}

func newNotify142Update(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.Notify142)

	switch data.Type {
	case 0xF:
		if data.P1 == 538 {
			return castingUpdate{
				streamID:  streamID,
				subjectID: uint64(b.SubjectID),

				castingInfo: nil,
			}
		}
	case 0x22:
		return lockonUpdate{
			streamID:  streamID,
			subjectID: uint64(b.SubjectID),

			lockonMarker: int(data.P1),
		}
	}
	return nil
}

type lockonUpdate struct {
	streamID  int
	subjectID uint64

	lockonMarker int
}

func (u lockonUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	return validateEntityUpdate(streams, u.streamID, u.subjectID, u.modifyFunc)
}

func (u lockonUpdate) modifyFunc(stream *models.Stream, entity *models.Entity) ([]models.StreamEvent, []models.EntityEvent, error) {
	entity.LockonMarker = u.lockonMarker

	return nil, []models.EntityEvent{models.EntityEvent{
		StreamID: u.streamID,
		EntityID: u.subjectID,
		Type: models.UpdateLockonMarker{
			LockonMarker: u.lockonMarker,
		},
	}}, nil
}
