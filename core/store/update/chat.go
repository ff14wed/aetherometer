package update

import (
	"fmt"
	"strings"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
)

const (
	ChatTypeParty         = 0x01
	ChatTypeLS            = 0x02
	ChatTypeFC            = 0x03
	ChatTypeNoviceNetwork = 0x04

	ChatZoneTypeSay   = 0x0A
	ChatZoneTypeYell  = 0x1E
	ChatZoneTypeShout = 0x0B
)

func init() {
	registerIngressHandler(new(datatypes.ChatFrom), newChatFromUpdate)
	registerEgressHandler(new(datatypes.ChatTo), newChatToUpdate)
	registerIngressHandler(new(datatypes.ChatFromXWorld), newChatFromXWorldUpdate)

	registerIngressHandler(new(datatypes.Chat), newChatUpdate)
	registerEgressHandler(new(datatypes.EgressChat), newEgressChatUpdate)

	registerIngressHandler(new(datatypes.FreeCompanyResult), newFCResultUpdate)

	registerIngressHandler(new(datatypes.ChatXWorld), newChatXWorldUpdate)
	registerEgressHandler(new(datatypes.EgressChatXWorld), newEgressChatXWorldUpdate)

	registerIngressHandler(new(datatypes.ChatZone), newChatZoneUpdate)
	registerEgressHandler(new(datatypes.EgressChatZone), newEgressChatZoneUpdate)
}

func newChatFromUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.ChatFrom)

	speakerWorld := d.WorldData.Lookup(int(data.WorldID))

	return chatUpdate{
		streamID: streamID,
		chatEvent: models.ChatEvent{
			ChannelType: "Private",

			ContentID: data.FromCharacterID,
			World:     speakerWorld,
			Name:      data.FromName.String(),

			Message: data.Message.String(),
		},
	}
}
func newChatToUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.ChatTo)

	toWorld := d.WorldData.Lookup(int(data.WorldID))

	return chatUpdate{
		streamID: streamID,
		chatEvent: models.ChatEvent{
			ChannelID:   data.ChannelID,
			ChannelType: "PrivateTo",

			ContentID: data.ToCharacterID,
			World:     toWorld,
			Name:      data.ToName.String(),

			Message: data.Message.String(),
		},
	}
}
func newChatFromXWorldUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.ChatFromXWorld)

	speakerWorld := d.WorldData.Lookup(int(data.WorldID))

	return chatUpdate{
		streamID: streamID,
		chatEvent: models.ChatEvent{
			ChannelType: "Private",

			ContentID: data.FromCharacterID,
			EntityID:  uint64(data.FromEntityID),
			World:     speakerWorld,
			Name:      data.FromName.String(),

			Message: data.Message.String(),
		},
	}
}

func channelTypeFromChannelID(channelID uint64) string {
	chatType := (channelID & 0xFF00000000) >> 32
	switch chatType {
	case ChatTypeParty:
		return "Party"
	case ChatTypeLS:
		return "Linkshell"
	case ChatTypeFC:
		return "FreeCompany"
	case ChatTypeNoviceNetwork:
		return "NoviceNetwork"
	}
	return fmt.Sprintf("Unknown_%d", chatType)
}

func channelWorldFromChannelID(channelID uint64) int {
	return int((channelID & 0x00FF000000000000) >> 48)
}

func newChatUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.Chat)

	channelType := channelTypeFromChannelID(data.ChannelID)
	channelWorld := d.WorldData.Lookup(channelWorldFromChannelID(data.ChannelID))
	speakerWorld := d.WorldData.Lookup(int(data.WorldID))

	return chatUpdate{
		streamID: streamID,
		chatEvent: models.ChatEvent{
			ChannelID:    data.ChannelID,
			ChannelWorld: channelWorld,
			ChannelType:  channelType,

			ContentID: data.SpeakerCharacterID,
			EntityID:  uint64(data.SpeakerEntityID),
			World:     speakerWorld,
			Name:      data.SpeakerName.String(),

			Message: data.Message.String(),
		},
	}
}

func newEgressChatUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.EgressChat)

	channelType := channelTypeFromChannelID(data.ChannelID)
	channelWorld := d.WorldData.Lookup(channelWorldFromChannelID(data.ChannelID))

	return chatUpdate{
		streamID: streamID,
		chatEvent: models.ChatEvent{
			ChannelID:    data.ChannelID,
			ChannelWorld: channelWorld,
			ChannelType:  channelType,

			Message: data.Message.String(),
		},
	}
}

func newFCResultUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.FreeCompanyResult)

	var message string

	if data.Type == 0xF {
		message = "logged in"
	} else if data.Type == 0x10 {
		message = "logged out"
	} else {
		message = fmt.Sprintf("Unknown_0x%x_0x%x_0x%x_0x%x", data.Type, data.Result, data.UpdateStatus, data.Identity)
	}

	return chatUpdate{
		streamID: streamID,
		chatEvent: models.ChatEvent{
			ChannelID:   data.FreeCompanyID,
			ChannelType: "FreeCompanyResult",

			ContentID: data.TargetCharacterID,
			Name:      data.FreeCompanyName.String(),

			Message: fmt.Sprintf("%s has %s.", data.TargetName.String(), message),
		},
	}
}

func newChatXWorldUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.ChatXWorld)

	// The assumption is this is only used for cross-world linkshell chat
	// since cross world party chat appears to be something else
	channelType := "CrossWorldLinkshell"
	speakerWorld := d.WorldData.Lookup(int(data.WorldID))

	return chatUpdate{
		streamID: streamID,
		chatEvent: models.ChatEvent{
			ChannelID:    data.ChannelID,
			ChannelWorld: new(models.World),
			ChannelType:  channelType,

			ContentID: data.SpeakerCharacterID,
			EntityID:  uint64(data.SpeakerEntityID),
			World:     speakerWorld,
			Name:      data.SpeakerName.String(),

			Message: data.Message.String(),
		},
	}
}

func newEgressChatXWorldUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.EgressChatXWorld)

	// The assumption is this is only used for cross-world linkshell chat
	// since cross world party chat appears to be something else
	channelType := "CrossWorldLinkshell"

	return chatUpdate{
		streamID: streamID,
		chatEvent: models.ChatEvent{
			ChannelID:    data.ChannelID,
			ChannelWorld: new(models.World),
			ChannelType:  channelType,

			Message: data.Message.String(),
		},
	}
}

func channelTypeFromChatZoneType(t uint16) string {
	if t == ChatZoneTypeSay {
		return "ZoneChatSay"
	} else if t == ChatZoneTypeShout {
		return "ZoneChatShout"
	} else if t == ChatZoneTypeYell {
		return "ZoneChatYell"
	}
	return fmt.Sprintf("ZoneChatUnknown%d", t)
}

func newChatZoneUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.ChatZone)

	speakerWorld := d.WorldData.Lookup(int(data.WorldID))
	channelType := channelTypeFromChatZoneType(data.Type)

	return chatUpdate{
		streamID: streamID,
		chatEvent: models.ChatEvent{
			ChannelType: channelType,

			ContentID: data.CharacterID,
			EntityID:  uint64(data.EntityID),
			World:     speakerWorld,
			Name:      data.SpeakerName.String(),

			Message: data.Message.String(),
		},
	}
}

func newEgressChatZoneUpdate(streamID int, b *xivnet.Block, d *datasheet.Collection) store.Update {
	data := b.Data.(*datatypes.EgressChatZone)

	channelType := channelTypeFromChatZoneType(data.Type)

	return chatUpdate{
		streamID: streamID,
		chatEvent: models.ChatEvent{
			ChannelType: channelType,

			Message: data.Message.String(),
		},
	}
}

type chatUpdate struct {
	streamID int

	chatEvent models.ChatEvent
}

func (u chatUpdate) ModifyStore(streams *store.Streams) ([]models.StreamEvent, []models.EntityEvent, error) {
	stream, found := streams.Map[u.streamID]
	if !found {
		return nil, nil, ErrorStreamNotFound
	}

	// If chat type is Zone, then the channel world is CurrentWorld.
	if strings.HasPrefix(u.chatEvent.ChannelType, "Zone") {
		u.chatEvent.ChannelWorld = &stream.CurrentWorld
	}

	// Using ContentID to check if it was an egress chat or not
	if u.chatEvent.ContentID == 0 {
		// If it was an egress chat then speaker is self.
		charName := "Me"
		if entity, ok := stream.EntitiesMap[stream.CharacterID]; ok {
			charName = entity.Name
		}
		u.chatEvent.EntityID = stream.CharacterID
		u.chatEvent.World = &stream.HomeWorld
		u.chatEvent.Name = charName
	}

	// If no ChannelWorld provided, then channel is assumed to be home world
	// unless it's cross-world
	if !strings.HasPrefix(u.chatEvent.ChannelType, "Cross") {
		if u.chatEvent.ChannelWorld == nil || u.chatEvent.ChannelWorld.Name == "" {
			u.chatEvent.ChannelWorld = &stream.HomeWorld
		}
	}

	if u.chatEvent.ChannelType == "FreeCompanyResult" {
		u.chatEvent.World = &stream.HomeWorld
	}

	return []models.StreamEvent{{
		StreamID: u.streamID,
		Type:     u.chatEvent,
	}}, nil, nil
}
