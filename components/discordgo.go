package components

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// InitDiscord initializes and caches the discord client
// reads the discord endpoint from the environment variable DISCORD_ENDPOINT, example: https://discordapp.com/
// reads the discord bot token from the environment variable DISCORD_BOT_TOKEN
func InitDiscord() (err error) {
	// create a new Discordgo Bot Client
	dhelpers.SetDiscordEndpoints(os.Getenv("DISCORD_ENDPOINT"))
	fmt.Println("set Discord Endpoint API URL to", discordgo.EndpointAPI)
	fmt.Println("creating Discord Client, Token length:", len(os.Getenv("DISCORD_BOT_TOKEN")))
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		return err
	}

	// cache client
	cache.SetDiscord(dg)

	return nil
}
