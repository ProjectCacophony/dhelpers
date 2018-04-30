package components

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// InitDiscord initializes and caches the discord client
// reads the discord endpoint from the environment variable DISCORD_ENDPOINT, example: https://discordapp.com/
// reads the discord bot token from the environment variable DISCORD_BOT_TOKEN
// this is only to open a gateway connection with a fixed bot token, for everything else cache.GetEDiscord should be used to get a customised session
func InitDiscord() (err error) {
	// create a new Discordgo Bot Client
	dhelpers.SetDiscordEndpoints(os.Getenv("DISCORD_ENDPOINT"))
	cache.GetLogger().Infoln("set Discord Endpoint API URL to", discordgo.EndpointAPI)
	cache.GetLogger().Infoln("creating Discord Client, Token length:", len(os.Getenv("DISCORD_BOT_TOKEN")))
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		return err
	}

	// cache client
	cache.SetDiscord(dg)

	return nil
}
