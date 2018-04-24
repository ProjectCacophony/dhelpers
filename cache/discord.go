package cache

import (
	"sync"

	"os"

	"github.com/Seklfreak/Robyul2/cache"
	"github.com/bwmarrin/discordgo"
)

var (
	discord       *discordgo.Session
	discordMutex  sync.RWMutex
	eventSessions = make(map[string]*discordgo.Session)
)

// SetDiscord caches a discordgo session for future use
func SetDiscord(s *discordgo.Session) {
	discordMutex.Lock()
	discord = s
	discordMutex.Unlock()
}

// GetDiscord returns a cached discordgo session
func GetDiscord() *discordgo.Session {
	discordMutex.RLock()
	defer discordMutex.RUnlock()

	return discord
}

// GetEDiscord gets or create a discord session for an Event
// reads the discord bot token from DISCORD_BOT_TOKEN_<bot user id>
func GetEDiscord(botID string) *discordgo.Session {
	discordMutex.Lock()
	defer discordMutex.Unlock()
	if _, ok := eventSessions[botID]; ok {
		return eventSessions[botID]
	}

	newSession, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN_"+botID))
	if err != nil {
		cache.GetLogger().WithField("module", "cache").Errorln("Error creating discord session for", botID+":", err.Error())
		return nil
	}
	eventSessions[botID] = newSession
	return newSession
}
