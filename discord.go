package dhelpers

import "gitlab.com/project-d-collab/dhelpers/cache"

// GoType starts a goroutine to start a typing indicator in a channel
func (event EventContainer) GoType(channelID string) {
	go func() {
		defer func() {
			recover() // nolint: errcheck, gas
		}()

		cache.GetEDiscord(event.BotUserID).ChannelTyping(channelID) // nolint: errcheck, gas
	}()
}
