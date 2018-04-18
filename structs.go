package dhelpers

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type EventContainer struct {
	Type           EventType
	Data           []byte
	ReceivedAt     time.Time
	GatewayStarted time.Time
}

type EventChannelCreate struct {
	Alias         string
	Prefix        string
	Event         *discordgo.ChannelCreate
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventChannelDelete struct {
	Alias         string
	Prefix        string
	Event         *discordgo.ChannelDelete
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventChannelPinsUpdate struct {
	Alias         string
	Prefix        string
	Event         *discordgo.ChannelPinsUpdate
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventChannelUpdate struct {
	Alias         string
	Prefix        string
	Event         *discordgo.ChannelUpdate
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventGuildBanAdd struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildBanAdd
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventGuildBanRemove struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildBanRemove
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventGuildCreate struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildCreate
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventGuildDelete struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildDelete
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventGuildEmojisUpdate struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildEmojisUpdate
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventGuildMemberAdd struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildMemberAdd
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventGuildMemberRemove struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildMemberRemove
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventGuildMemberUpdate struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildMemberUpdate
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventGuildMembersChunk struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildMembersChunk
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventGuildRoleCreate struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildRoleCreate
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}
type EventGuildRoleDelete struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildRoleDelete
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventGuildRoleUpdate struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildRoleUpdate
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventGuildUpdate struct {
	Alias         string
	Prefix        string
	Event         *discordgo.GuildUpdate
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventMessageCreate struct {
	Alias         string
	Prefix        string
	Event         *discordgo.MessageCreate
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
	Args          []string
}

type EventMessageDelete struct {
	Alias         string
	Prefix        string
	Event         *discordgo.MessageDelete
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
	Args          []string
}

type EventMessageReactionAdd struct {
	Alias         string
	Prefix        string
	Event         *discordgo.MessageReactionAdd
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventMessageReactionRemove struct {
	Alias         string
	Prefix        string
	Event         *discordgo.MessageReactionRemove
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventMessageReactionRemoveAll struct {
	Alias         string
	Prefix        string
	Event         *discordgo.MessageReactionRemoveAll
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}

type EventMessageUpdate struct {
	Alias         string
	Prefix        string
	Event         *discordgo.MessageUpdate
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
	Args          []string
}

type EventPresenceUpdate struct {
	Alias         string
	Prefix        string
	Event         *discordgo.PresenceUpdate
	BotUser       *discordgo.User
	SourceChannel *discordgo.Channel
	SourceGuild   *discordgo.Guild
}
