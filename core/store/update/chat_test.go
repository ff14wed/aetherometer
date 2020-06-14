package update_test

import (
	"fmt"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Chat Update", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streams   *store.Streams
		d         *datasheet.Collection
		streamID  int
		subjectID uint64
		generator update.Generator

		expectedChatEvent models.ChatEvent
	)

	BeforeEach(func() {
		*testEnv = genericSetup()
		b = testEnv.b
		streams = testEnv.streams
		d = testEnv.d
		streamID = testEnv.streamID
		subjectID = testEnv.subjectID
		generator = testEnv.generator

		d.WorldData = datasheet.WorldStore{
			123: {Key: 123, Name: "Foo"},
			456: {Key: 123, Name: "Bar"},
		}
	})

	Describe("ChatFrom", func() {
		BeforeEach(func() {
			chatData := &datatypes.ChatFrom{
				FromCharacterID: 0x004000170000001,
				WorldID:         123,
				FromName:        datatypes.StringToEntityName("Sender"),
				Message:         datatypes.StringToChatMessage("Private message"),
			}

			b.Data = chatData

			expectedChatEvent = models.ChatEvent{
				ChannelID:    0x0,
				ChannelWorld: models.World{ID: 456, Name: "Bar"},
				ChannelType:  "Private",

				ContentID: 0x004000170000001,
				EntityID:  0x0,
				World:     models.World{ID: 123, Name: "Foo"},
				Name:      "Sender",
				Message:   "Private message",
			}
		})

		It("generates a StreamEvent for the chat event", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(entityEvents).To(BeEmpty())

			Expect(streamEvents).To(ConsistOf(models.StreamEvent{
				StreamID: streamID,
				Type:     expectedChatEvent,
			}))
		})

		streamValidationTests(testEnv, false)
	})

	Describe("ChatTo", func() {
		BeforeEach(func() {
			chatData := &datatypes.ChatTo{
				ChannelID:     0xABCD,
				ToCharacterID: 0x004000170000002,
				WorldID:       123,
				ToName:        datatypes.StringToEntityName("Recipient"),
				Message:       datatypes.StringToChatMessage("Private message"),
			}

			b.Data = chatData

			expectedChatEvent = models.ChatEvent{
				ChannelID:    0xABCD,
				ChannelWorld: models.World{ID: 456, Name: "Bar"},
				ChannelType:  "PrivateTo",

				ContentID: 0x004000170000002,
				EntityID:  0x0,
				World:     models.World{ID: 123, Name: "Foo"},
				Name:      "Recipient",
				Message:   "Private message",
			}
		})

		It("generates a StreamEvent for the chat event", func() {
			u := generator.Generate(streamID, true, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(entityEvents).To(BeEmpty())

			Expect(streamEvents).To(ConsistOf(models.StreamEvent{
				StreamID: streamID,
				Type:     expectedChatEvent,
			}))
		})

		streamValidationTests(testEnv, true)
	})

	Describe("Chat", func() {
		chatTypes := []uint64{
			update.ChatTypeParty,
			update.ChatTypeLS,
			update.ChatTypeFC,
			update.ChatTypeNoviceNetwork,
		}
		channelTypes := []string{
			"Party",
			"Linkshell",
			"FreeCompany",
			"NoviceNetwork",
		}

		for i, t := range channelTypes {
			channelType := t
			chatType := chatTypes[i]
			Context(fmt.Sprintf("with ChannelType %s", channelType), func() {
				BeforeEach(func() {
					var channelID uint64 = 0x007B00000000ABCD
					mask := chatType << 32
					channelID = channelID | mask

					chatData := &datatypes.Chat{
						SpeakerCharacterID: 0x004000170000001,
						SpeakerEntityID:    0x12345678,
						WorldID:            456,
						SpeakerName:        datatypes.StringToEntityName("Sender"),
						Message:            datatypes.StringToChatMessage("Blah blah"),
					}

					chatData.ChannelID = channelID

					b.Data = chatData

					expectedChatEvent = models.ChatEvent{
						ChannelID:    channelID,
						ChannelWorld: models.World{ID: 123, Name: "Foo"},
						ChannelType:  channelType,

						ContentID: 0x004000170000001,
						EntityID:  0x12345678,
						World:     models.World{ID: 456, Name: "Bar"},
						Name:      "Sender",
						Message:   "Blah blah",
					}
				})

				It("generates a StreamEvent for the chat event", func() {
					u := generator.Generate(streamID, false, b)
					Expect(u).ToNot(BeNil())
					streamEvents, entityEvents, err := u.ModifyStore(streams)
					Expect(err).ToNot(HaveOccurred())
					Expect(entityEvents).To(BeEmpty())

					Expect(streamEvents).To(ConsistOf(models.StreamEvent{
						StreamID: streamID,
						Type:     expectedChatEvent,
					}))
				})
				streamValidationTests(testEnv, false)
			})
		}
	})

	Describe("EgressChat", func() {
		chatTypes := []uint64{
			update.ChatTypeParty,
			update.ChatTypeLS,
			update.ChatTypeFC,
			update.ChatTypeNoviceNetwork,
		}
		channelTypes := []string{
			"Party",
			"Linkshell",
			"FreeCompany",
			"NoviceNetwork",
		}

		for i, t := range channelTypes {
			channelType := t
			chatType := chatTypes[i]

			Context(fmt.Sprintf("with ChannelType %s", channelType), func() {
				BeforeEach(func() {
					var channelID uint64 = 0x007B00000000ABCD
					mask := chatType << 32
					channelID = channelID | mask

					chatData := &datatypes.EgressChat{
						Message: datatypes.StringToChatMessage("Blah blah"),
					}

					chatData.ChannelID = channelID

					b.Data = chatData

					expectedChatEvent = models.ChatEvent{
						ChannelID:    channelID,
						ChannelWorld: models.World{ID: 123, Name: "Foo"},
						ChannelType:  channelType,

						ContentID: 0,
						EntityID:  subjectID,
						World:     models.World{ID: 456, Name: "Bar"},
						Name:      "Test Subject",
						Message:   "Blah blah",
					}
				})

				It("fills in the sender information and generates a StreamEvent for the chat event", func() {
					u := generator.Generate(streamID, true, b)
					Expect(u).ToNot(BeNil())
					streamEvents, entityEvents, err := u.ModifyStore(streams)
					Expect(err).ToNot(HaveOccurred())
					Expect(entityEvents).To(BeEmpty())

					Expect(streamEvents).To(ConsistOf(models.StreamEvent{
						StreamID: streamID,
						Type:     expectedChatEvent,
					}))
				})
				streamValidationTests(testEnv, true)
			})
		}
	})

	Describe("FCResult", func() {
		Context("when a user has logged in", func() {
			BeforeEach(func() {
				chatData := &datatypes.FreeCompanyResult{
					FreeCompanyID:     0xABCD,
					TargetCharacterID: 0x004000170000001,
					Type:              0xF,
					FreeCompanyName:   datatypes.StringToFCName("Free the Company"),
					TargetName:        datatypes.StringToEntityName("Some Employee"),
				}

				b.Data = chatData

				expectedChatEvent = models.ChatEvent{
					ChannelID:    0xABCD,
					ChannelWorld: models.World{ID: 456, Name: "Bar"},
					ChannelType:  "FreeCompanyResult",

					ContentID: 0x004000170000001,
					Name:      "Free the Company",
					Message:   "Some Employee has logged in.",
				}
			})

			It("generates a StreamEvent for the chat event", func() {
				u := generator.Generate(streamID, false, b)
				Expect(u).ToNot(BeNil())
				streamEvents, entityEvents, err := u.ModifyStore(streams)
				Expect(err).ToNot(HaveOccurred())
				Expect(entityEvents).To(BeEmpty())

				Expect(streamEvents).To(ConsistOf(models.StreamEvent{
					StreamID: streamID,
					Type:     expectedChatEvent,
				}))
			})

			streamValidationTests(testEnv, false)
		})

		Context("when a user has logged out", func() {
			BeforeEach(func() {
				chatData := &datatypes.FreeCompanyResult{
					FreeCompanyID:     0xABCD,
					TargetCharacterID: 0x004000170000001,
					Type:              0x10,
					FreeCompanyName:   datatypes.StringToFCName("Free the Company"),
					TargetName:        datatypes.StringToEntityName("Some Employee"),
				}

				b.Data = chatData

				expectedChatEvent = models.ChatEvent{
					ChannelID:    0xABCD,
					ChannelWorld: models.World{ID: 456, Name: "Bar"},
					ChannelType:  "FreeCompanyResult",

					ContentID: 0x004000170000001,
					Name:      "Free the Company",
					Message:   "Some Employee has logged out.",
				}
			})

			It("generates a StreamEvent for the chat event", func() {
				u := generator.Generate(streamID, false, b)
				Expect(u).ToNot(BeNil())
				streamEvents, entityEvents, err := u.ModifyStore(streams)
				Expect(err).ToNot(HaveOccurred())
				Expect(entityEvents).To(BeEmpty())

				Expect(streamEvents).To(ConsistOf(models.StreamEvent{
					StreamID: streamID,
					Type:     expectedChatEvent,
				}))
			})

			streamValidationTests(testEnv, false)
		})
	})

	Describe("ChatXWorld", func() {
		BeforeEach(func() {
			var channelID uint64 = 0x007B00060000ABCD

			chatData := &datatypes.ChatXWorld{
				SpeakerCharacterID: 0x004000170000001,
				SpeakerEntityID:    0x12345678,
				WorldID:            456,
				SpeakerName:        datatypes.StringToEntityName("Sender"),
				Message:            datatypes.StringToChatMessage("Blah blah"),
			}

			chatData.ChannelID = channelID

			b.Data = chatData

			expectedChatEvent = models.ChatEvent{
				ChannelID:   channelID,
				ChannelType: "CrossWorldLinkshell",

				ContentID: 0x004000170000001,
				EntityID:  0x12345678,
				World:     models.World{ID: 456, Name: "Bar"},
				Name:      "Sender",
				Message:   "Blah blah",
			}
		})

		It("generates a StreamEvent for the chat event", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(entityEvents).To(BeEmpty())

			Expect(streamEvents).To(ConsistOf(models.StreamEvent{
				StreamID: streamID,
				Type:     expectedChatEvent,
			}))
		})
		streamValidationTests(testEnv, false)
	})

	Describe("EgressChatXWorld", func() {
		BeforeEach(func() {
			var channelID uint64 = 0x007B00060000ABCD

			chatData := &datatypes.EgressChatXWorld{
				Message: datatypes.StringToChatMessage("Blah blah"),
			}

			chatData.ChannelID = channelID

			b.Data = chatData

			expectedChatEvent = models.ChatEvent{
				ChannelID:   channelID,
				ChannelType: "CrossWorldLinkshell",

				ContentID: 0,
				EntityID:  subjectID,
				World:     models.World{ID: 456, Name: "Bar"},
				Name:      "Test Subject",
				Message:   "Blah blah",
			}
		})

		It("fills in the sender information and generates a StreamEvent for the chat event", func() {
			u := generator.Generate(streamID, true, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(entityEvents).To(BeEmpty())

			Expect(streamEvents).To(ConsistOf(models.StreamEvent{
				StreamID: streamID,
				Type:     expectedChatEvent,
			}))
		})
		streamValidationTests(testEnv, true)
	})

	Describe("ChatZone", func() {
		chatTypes := []uint16{
			update.ChatZoneTypeSay,
			update.ChatZoneTypeShout,
			update.ChatZoneTypeYell,
		}
		channelTypes := []string{
			"ZoneChatSay",
			"ZoneChatShout",
			"ZoneChatYell",
		}

		for i, t := range channelTypes {
			channelType := t
			chatType := chatTypes[i]
			Context(fmt.Sprintf("with ChannelType %s", channelType), func() {
				BeforeEach(func() {

					chatData := &datatypes.ChatZone{
						CharacterID: 0x004000170000001,
						EntityID:    0x12345678,
						WorldID:     456,
						Type:        chatType,
						SpeakerName: datatypes.StringToEntityName("Sender"),
						Message:     datatypes.StringToChatMessage("Blah blah"),
					}

					b.Data = chatData

					expectedChatEvent = models.ChatEvent{
						ChannelWorld: models.World{ID: 123, Name: "Foo"},
						ChannelType:  channelType,

						ContentID: 0x004000170000001,
						EntityID:  0x12345678,
						World:     models.World{ID: 456, Name: "Bar"},
						Name:      "Sender",
						Message:   "Blah blah",
					}
				})

				It("generates a StreamEvent for the chat event", func() {
					u := generator.Generate(streamID, false, b)
					Expect(u).ToNot(BeNil())
					streamEvents, entityEvents, err := u.ModifyStore(streams)
					Expect(err).ToNot(HaveOccurred())
					Expect(entityEvents).To(BeEmpty())

					Expect(streamEvents).To(ConsistOf(models.StreamEvent{
						StreamID: streamID,
						Type:     expectedChatEvent,
					}))
				})
				streamValidationTests(testEnv, false)
			})
		}
	})

	Describe("EgressChat", func() {
		chatTypes := []uint16{
			update.ChatZoneTypeSay,
			update.ChatZoneTypeShout,
			update.ChatZoneTypeYell,
		}
		channelTypes := []string{
			"ZoneChatSay",
			"ZoneChatShout",
			"ZoneChatYell",
		}

		for i, t := range channelTypes {
			channelType := t
			chatType := chatTypes[i]

			Context(fmt.Sprintf("with ChannelType %s", channelType), func() {
				BeforeEach(func() {
					chatData := &datatypes.EgressChatZone{
						Type:    chatType,
						Message: datatypes.StringToChatMessage("Blah blah"),
					}

					b.Data = chatData

					expectedChatEvent = models.ChatEvent{
						ChannelWorld: models.World{ID: 123, Name: "Foo"},
						ChannelType:  channelType,

						ContentID: 0,
						EntityID:  subjectID,
						World:     models.World{ID: 456, Name: "Bar"},
						Name:      "Test Subject",
						Message:   "Blah blah",
					}
				})

				It("fills in the sender information and generates a StreamEvent for the chat event", func() {
					u := generator.Generate(streamID, true, b)
					Expect(u).ToNot(BeNil())
					streamEvents, entityEvents, err := u.ModifyStore(streams)
					Expect(err).ToNot(HaveOccurred())
					Expect(entityEvents).To(BeEmpty())

					Expect(streamEvents).To(ConsistOf(models.StreamEvent{
						StreamID: streamID,
						Type:     expectedChatEvent,
					}))
				})
				streamValidationTests(testEnv, true)
			})
		}
	})
})
