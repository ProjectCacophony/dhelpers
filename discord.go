package dhelpers

import (
	"time"

	"strings"

	"gitlab.com/Cacophony/dhelpers/cache"
)

// GoType starts a goroutine to start a typing indicator in a channel
// it will use the channelID if given, if not it will try to use the channelID from the event,
// the following events are supported at the moment:
// MESSAGE_CREATE, MESSAGE_UPDATE, MESSAGE_DELETE, MESSAGE_REACTION_ADD, MESSAGE_REACTION_REMOVE,
//  MESSAGE_REACTION_REMOVE_ALL, CHANNEL_CREATE, CHANNEL_UPDATE, CHANNEL_DELETE, CHANNEL_PINS_UPDATE
func (event EventContainer) GoType(channelID ...string) {
	go func() {
		defer func() {
			recover() // nolint: errcheck, gas
		}()

		var typingChannel string

		if len(channelID) > 0 {
			// use channel ID if given
			typingChannel = channelID[0]
		} else {
			// or try to use channel from event
			switch event.Type {
			case MessageCreateEventType:
				typingChannel = event.MessageCreate.ChannelID
			case MessageUpdateEventType:
				typingChannel = event.MessageUpdate.ChannelID
			case MessageDeleteEventType:
				typingChannel = event.MessageDelete.ChannelID
			case MessageReactionAddEventType:
				typingChannel = event.MessageReactionAdd.ChannelID
			case MessageReactionRemoveEventType:
				typingChannel = event.MessageReactionRemove.ChannelID
			case MessageReactionRemoveAllEventType:
				typingChannel = event.MessageReactionRemoveAll.ChannelID
			case ChannelCreateEventType:
				typingChannel = event.ChannelCreate.ID
			case ChannelUpdateEventType:
				typingChannel = event.ChannelUpdate.ID
			case ChannelDeleteEventType:
				typingChannel = event.ChannelDelete.ID
			case ChannelPinsUpdateEventType:
				typingChannel = event.ChannelPinsUpdate.ChannelID
			}
		}

		if typingChannel != "" {
			cache.GetEDiscord(event.BotUserID).ChannelTyping(typingChannel) // nolint: errcheck, gas
		}
	}()
}

// DiscordTime returns a time formatted to be used in Embeds
func DiscordTime(theTime time.Time) string {
	return theTime.Format(time.RFC3339)
}

// CleanURL makes a URL posted in discord ready to use for further usage
func CleanURL(uncleanedURL string) (url string) {
	if strings.HasPrefix(uncleanedURL, "<") {
		uncleanedURL = strings.TrimLeft(uncleanedURL, "<")
	}
	if strings.HasSuffix(uncleanedURL, ">") {
		uncleanedURL = strings.TrimRight(uncleanedURL, ">")
	}
	return uncleanedURL
}
