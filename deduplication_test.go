package dhelpers

import (
	"testing"

	"os"

	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gitlab.com/Cacophony/dhelpers/cache"
)

func init() {
	// init logger
	cache.SetLogger(logrus.NewEntry(logrus.New()))
	// init redis
	cache.SetRedisClient(redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "",
		DB:       0,
	}))
}

func TestGetEventKey(t *testing.T) {
	v := GetEventKey(&discordgo.GuildCreate{})
	if v != "cacophony:gateway:event-GUILD_CREATE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-GUILD_CREATE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.GuildUpdate{})
	if v != "cacophony:gateway:event-GUILD_UPDATE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-GUILD_UPDATE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.GuildDelete{})
	if v != "cacophony:gateway:event-GUILD_DELETE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-GUILD_DELETE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.GuildMemberAdd{})
	if v != "cacophony:gateway:event-GUILD_MEMBER_ADD-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-GUILD_MEMBER_ADD-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.GuildMemberUpdate{})
	if v != "cacophony:gateway:event-GUILD_MEMBER_UPDATE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-GUILD_MEMBER_UPDATE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.GuildMemberRemove{})
	if v != "cacophony:gateway:event-GUILD_MEMBER_REMOVE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-GUILD_MEMBER_REMOVE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.GuildMembersChunk{})
	if v != "cacophony:gateway:event-GUILD_MEMBERS_CHUNK-5f1cb4e5159145f28dc6b9176b2c2ef4" {
		t.Error("Expected cacophony:gateway:event-GUILD_MEMBERS_CHUNK-5f1cb4e5159145f28dc6b9176b2c2ef4, got ", v)
	}
	v = GetEventKey(&discordgo.GuildRoleCreate{})
	if v != "cacophony:gateway:event-GUILD_ROLE_CREATE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-GUILD_ROLE_CREATE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.GuildRoleUpdate{})
	if v != "cacophony:gateway:event-GUILD_ROLE_UPDATE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-GUILD_ROLE_UPDATE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.GuildRoleDelete{})
	if v != "cacophony:gateway:event-GUILD_ROLE_DELETE-7215ee9c7d9dc229d2921a40e899ec5f" {
		t.Error("Expected cacophony:gateway:event-GUILD_ROLE_DELETE-7215ee9c7d9dc229d2921a40e899ec5f, got ", v)
	}
	v = GetEventKey(&discordgo.GuildEmojisUpdate{})
	if v != "cacophony:gateway:event-GUILD_EMOJIS_UPDATE-5f1cb4e5159145f28dc6b9176b2c2ef4" {
		t.Error("Expected ecacophony:gateway:event-GUILD_EMOJIS_UPDATE-5f1cb4e5159145f28dc6b9176b2c2ef4, got ", v)
	}
	v = GetEventKey(&discordgo.ChannelCreate{})
	if v != "cacophony:gateway:event-CHANNEL_CREATE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-CHANNEL_CREATE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.ChannelUpdate{})
	if v != "cacophony:gateway:event-CHANNEL_UPDATE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-CHANNEL_UPDATE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.ChannelDelete{})
	if v != "cacophony:gateway:event-CHANNEL_DELETE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-CHANNEL_DELETE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.MessageCreate{})
	if v != "cacophony:gateway:event-MESSAGE_CREATE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-MESSAGE_CREATE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.MessageUpdate{})
	if v != "cacophony:gateway:event-MESSAGE_UPDATE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-MESSAGE_UPDATE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.MessageDelete{})
	if v != "cacophony:gateway:event-MESSAGE_DELETE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-MESSAGE_DELETE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.PresenceUpdate{})
	if v != "cacophony:gateway:event-PRESENCE_UPDATE-d69245414ba43f1834c3ddc97701e130" {
		t.Error("Expected cacophony:gateway:event-PRESENCE_UPDATE-d69245414ba43f1834c3ddc97701e130, got ", v)
	}
	v = GetEventKey(&discordgo.ChannelPinsUpdate{})
	if v != "cacophony:gateway:event-CHANNEL_PINS_UPDATE-7215ee9c7d9dc229d2921a40e899ec5f" {
		t.Error("Expected cacophony:gateway:event-CHANNEL_PINS_UPDATE-7215ee9c7d9dc229d2921a40e899ec5f, got ", v)
	}
	v = GetEventKey(&discordgo.GuildBanAdd{})
	if v != "cacophony:gateway:event-GUILD_BAN_ADD-3b1cbb21fe5256cdfb5d6d1d659c56ff" {
		t.Error("Expected cacophony:gateway:event-GUILD_BAN_ADD-3b1cbb21fe5256cdfb5d6d1d659c56ff, got ", v)
	}
	v = GetEventKey(&discordgo.GuildBanRemove{})
	if v != "cacophony:gateway:event-GUILD_BAN_REMOVE-3b1cbb21fe5256cdfb5d6d1d659c56ff" {
		t.Error("Expected cacophony:gateway:event-GUILD_BAN_REMOVE-3b1cbb21fe5256cdfb5d6d1d659c56ff, got ", v)
	}
	v = GetEventKey(&discordgo.MessageReactionAdd{})
	if v != "cacophony:gateway:event-MESSAGE_REACTION_ADD-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-MESSAGE_REACTION_ADD-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.MessageReactionRemove{})
	if v != "cacophony:gateway:event-MESSAGE_REACTION_REMOVE-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-MESSAGE_REACTION_REMOVE-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(&discordgo.MessageReactionRemoveAll{})
	if v != "cacophony:gateway:event-MESSAGE_REACTION_REMOVE_ALL-60046f14c917c18a9a0f923e191ba0dc" {
		t.Error("Expected cacophony:gateway:event-MESSAGE_REACTION_REMOVE_ALL-60046f14c917c18a9a0f923e191ba0dc, got ", v)
	}
	v = GetEventKey(nil)
	if v != "" {
		t.Error("Expected , got ", v)
	}
}

func TestIsNewEvent(t *testing.T) {
	key1 := "cacophony:gateway:event-MESSAGE_CREATE-" + strconv.FormatInt(time.Now().Unix(), 10)
	key2 := "cacophony:gateway:event-MESSAGE_UPDATE-" + strconv.FormatInt(time.Now().Unix(), 10)
	v := IsNewEvent(cache.GetRedisClient(), "testing", key1)
	if !v {
		t.Error("Expected true, got ", v)
	}
	v = IsNewEvent(cache.GetRedisClient(), "testing", key1)
	if v {
		t.Error("Expected false, got ", v)
	}
	v = IsNewEvent(cache.GetRedisClient(), "testing", key2)
	if !v {
		t.Error("Expected true, got ", v)
	}
	v = IsNewEvent(redis.NewClient(&redis.Options{
		Addr:     "example.org",
		Password: "",
		DB:       0,
	}), "testing", key1)
	if v {
		t.Error("Expected false, got ", v)
	}
}
