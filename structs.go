package dhelpers

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Event struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             interface{}
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventChannelCreate struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.ChannelCreate
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventChannelDelete struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.ChannelDelete
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventPinsUpdate struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.ChannelPinsUpdate
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventChannelUpdate struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.ChannelUpdate
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventGuildBanAdd struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildBanAdd
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventGuildBanRemove struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildBanRemove
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventGuildCreate struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildCreate
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventGuildDelete struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildDelete
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventEmojisUpdate struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildEmojisUpdate
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventMemberAdd struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildMemberAdd
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventMemberRemove struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildMemberRemove
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventMemberUpdate struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildMemberUpdate
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventMembersChunk struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildMembersChunk
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventRoleCreate struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildRoleCreate
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}
type EventRoleDelete struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildRoleDelete
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventRoleUpdate struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildRoleUpdate
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventGuildUpdate struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.GuildUpdate
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventMessageCreate struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.MessageCreate
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventMessageDelete struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.MessageDelete
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventMessageReactionAdd struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.MessageReactionAdd
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventMessageReactionRemove struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.MessageReactionRemove
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventMessageReactionRemoveAll struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.MessageReactionRemoveAll
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventMessageUpdate struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.MessageUpdate
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}

type EventPresenceUpdate struct {
	Alias             string
	Type              EventType
	Prefix            string
	Event             discordgo.PresenceUpdate
	BotUser           *discordgo.User
	SourceChannel     *discordgo.Channel
	SourceGuild       *discordgo.Guild
	GatewayReceivedAt time.Time
	GatewayStarted    time.Time
}
