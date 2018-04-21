package cache

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	discord      *discordgo.Session
	discordMutex sync.RWMutex
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
