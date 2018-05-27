package dhelpers

import (
	"os"
	"testing"
	"time"

	"reflect"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/dhelpers/cache"
)

func TestCreateEventContainer(t *testing.T) {
	botUserID := os.Getenv("TESTING_DISCORD_BOTID")

	now := time.Now()
	example := EventContainer{
		ReceivedAt:     now,
		GatewayStarted: now,
		Key:            "foo-bar",
		BotUserID:      os.Getenv("TESTING_DISCORD_BOTID"),
	}
	v := CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", nil)
	if !reflect.DeepEqual(example, v) {
		t.Errorf("Expected %+v, got %+v", example, v)
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildCreate{})
	if v.GuildCreate == nil {
		t.Errorf("Expected GuildCreate to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildUpdate{})
	if v.GuildUpdate == nil {
		t.Errorf("Expected GuildUpdate to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildDelete{})
	if v.GuildDelete == nil {
		t.Errorf("Expected GuildDelete to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildMemberAdd{})
	if v.GuildMemberAdd == nil {
		t.Errorf("Expected GuildMemberAdd to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildMemberUpdate{})
	if v.GuildMemberUpdate == nil {
		t.Errorf("Expected GuildMemberUpdate to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildMemberRemove{})
	if v.GuildMemberRemove == nil {
		t.Errorf("Expected GuildMemberRemove to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildMembersChunk{})
	if v.GuildMembersChunk == nil {
		t.Errorf("Expected GuildMembersChunk to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildRoleCreate{})
	if v.GuildRoleCreate == nil {
		t.Errorf("Expected GuildRoleCreate to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildRoleUpdate{})
	if v.GuildRoleUpdate == nil {
		t.Errorf("Expected GuildRoleUpdate to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildRoleDelete{})
	if v.GuildRoleDelete == nil {
		t.Errorf("Expected GuildRoleUpdate to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildEmojisUpdate{})
	if v.GuildEmojisUpdate == nil {
		t.Errorf("Expected GuildEmojisUpdate to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.ChannelCreate{})
	if v.ChannelCreate == nil {
		t.Errorf("Expected ChannelCreate to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.ChannelUpdate{})
	if v.ChannelUpdate == nil {
		t.Errorf("Expected ChannelUpdate to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.ChannelDelete{})
	if v.ChannelDelete == nil {
		t.Errorf("Expected ChannelDelete to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.MessageCreate{
		Message: &discordgo.Message{
			GuildID: os.Getenv("TESTING_DISCORD_GUILDID"),
			Content: "test",
		},
	})
	if v.MessageCreate == nil {
		t.Errorf("Expected MessageCreate to not be nil")
	}
	if v.Prefix != "" {
		t.Error("Ewxpected Prefix to be , got ", v.Prefix)
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.MessageCreate{
		Message: &discordgo.Message{
			GuildID: os.Getenv("TESTING_DISCORD_GUILDID"),
			Content: "!test",
		},
	})
	if v.MessageCreate == nil {
		t.Errorf("Expected MessageCreate to not be nil")
	}
	if v.Prefix != "!" {
		t.Error("Expected Prefix to be !, got ", v.Prefix)
	}
	if len(v.Args) != 1 {
		t.Error("Expected 1 Args, got ", len(v.Args))
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.MessageUpdate{
		Message: &discordgo.Message{
			GuildID: os.Getenv("TESTING_DISCORD_GUILDID"),
			Content: "test",
		},
	})
	if v.MessageUpdate == nil {
		t.Errorf("Expected MessageUpdate to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.MessageDelete{})
	if v.MessageDelete == nil {
		t.Errorf("Expected MessageDelete to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.ChannelPinsUpdate{})
	if v.ChannelPinsUpdate == nil {
		t.Errorf("Expected ChannelPinsUpdate to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildBanAdd{})
	if v.GuildBanAdd == nil {
		t.Errorf("Expected GuildBanAdd to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.GuildBanRemove{})
	if v.GuildBanRemove == nil {
		t.Errorf("Expected GuildBanRemove to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.MessageReactionAdd{})
	if v.MessageReactionAdd == nil {
		t.Errorf("Expected MessageReactionAdd to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.MessageReactionRemove{})
	if v.MessageReactionRemove == nil {
		t.Errorf("Expected MessageReactionRemove to not be nil")
	}
	v = CreateEventContainer(now, now, cache.GetEDiscordGateway(botUserID), "foo-bar", &discordgo.MessageReactionRemoveAll{})
	if v.MessageReactionRemoveAll == nil {
		t.Errorf("Expected MessageReactionRemoveAll to not be nil")
	}
}
