package dhelpers

import (
	"time"

	"gitlab.com/Cacophony/dhelpers/cache"
)

// GoType starts a goroutine to start a typing indicator in a channel
func (event EventContainer) GoType(channelID string) {
	go func() {
		defer func() {
			recover() // nolint: errcheck, gas
		}()

		cache.GetEDiscord(event.BotUserID).ChannelTyping(channelID) // nolint: errcheck, gas
	}()
}

// DiscordTime returns a time formatted to be used in Embeds
func DiscordTime(theTime time.Time) string {
	return theTime.Format(time.RFC3339)
}
