package dhelpers

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestSetDiscordEndpoints(t *testing.T) {
	SetDiscordEndpoints("https://example.org/")
	if discordgo.EndpointDiscord != "https://example.org/" {
		t.Error("Expected https://example.org/, got ", discordgo.EndpointDiscord)
	}
	if discordgo.EndpointAPI != "https://example.org/"+"api/v"+discordgo.APIVersion+"/" {
		t.Error("Expected https://example.org/"+"api/v"+discordgo.APIVersion+"/, got ", discordgo.EndpointDiscord)
	}
	if discordgo.EndpointGuilds != "https://example.org/api/v"+discordgo.APIVersion+"/guilds/" {
		t.Error("Expected https://example.org/api/v"+discordgo.APIVersion+"/guilds/, got ", discordgo.EndpointDiscord)
	}
	if discordgo.EndpointChannels != "https://example.org/api/v"+discordgo.APIVersion+"/channels/" {
		t.Error("Expected https://example.org/api/v"+discordgo.APIVersion+"/channels/, got ", discordgo.EndpointDiscord)
	}
	if discordgo.EndpointUsers != "https://example.org/api/v"+discordgo.APIVersion+"/users/" {
		t.Error("Expected https://example.org/api/v"+discordgo.APIVersion+"/users/, got ", discordgo.EndpointDiscord)
	}
	if discordgo.EndpointGateway != "https://example.org/api/v"+discordgo.APIVersion+"/gateway" {
		t.Error("Expected https://example.org/api/v"+discordgo.APIVersion+"/gateway, got ", discordgo.EndpointDiscord)
	}
	if discordgo.EndpointGatewayBot != "https://example.org/api/v"+discordgo.APIVersion+"/gateway/bot" {
		t.Error("Expected https://example.org/api/v"+discordgo.APIVersion+"/gateway/bot, got ", discordgo.EndpointDiscord)
	}
	if discordgo.EndpointWebhooks != "https://example.org/api/v"+discordgo.APIVersion+"/webhooks/" {
		t.Error("Expected https://example.org/api/v"+discordgo.APIVersion+"/webhooks/, got ", discordgo.EndpointDiscord)
	}
}
