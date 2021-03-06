package dhelpers

import (
	"testing"

	"time"

	"os"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/dhelpers/cache"
)

func TestGoType(t *testing.T) {
	channelID := os.Getenv("TESTING_DISCORD_CHANNELID")
	event := EventContainer{
		BotUserID: os.Getenv("TESTING_DISCORD_BOTID"),
	}

	received := make(chan bool)

	cache.GetEDiscordGateway(event.BotUserID).AddHandler(func(s *discordgo.Session, m *discordgo.TypingStart) {
		if m.ChannelID != channelID {
			return
		}
		if m.UserID != event.BotUserID {
			return
		}

		received <- true
	})

	event.GoType(channelID)

	select {
	case res := <-received:
		if !res {
			t.Error("Did not receive successful gateway event")
		}
	case <-time.After(5 * time.Second):
		t.Error("Gateway event timed out")
	}
}

func TestDiscordTime(t *testing.T) {
	now := time.Now()
	v := DiscordTime(now)
	if v != now.Format(time.RFC3339) {
		t.Error("Expected ", now.Format(time.RFC3339), ", got ", v)
	}
}

func TestCleanURL(t *testing.T) {
	v := CleanURL("https://example.org")
	if v != "https://example.org" {
		t.Error("Expected https://example.org, got ", v)
	}
	v = CleanURL("<https://example.org>")
	if v != "https://example.org" {
		t.Error("Expected https://example.org, got ", v)
	}
}
