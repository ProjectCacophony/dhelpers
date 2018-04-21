package dhelpers

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// EventContainer is a container for all events sent to Lambdas or the SQS Queue
type EventContainer struct {
	Type           EventType
	Key            string
	ReceivedAt     time.Time
	GatewayStarted time.Time
	Modules        []string
	Destinations   []DestinationData
	Prefix         string
	Args           []string
	BotUserID      string
	// Events
	ChannelCreate            *discordgo.ChannelCreate
	ChannelDelete            *discordgo.ChannelDelete
	ChannelPinsUpdate        *discordgo.ChannelPinsUpdate
	ChannelUpdate            *discordgo.ChannelUpdate
	GuildBanAdd              *discordgo.GuildBanAdd
	GuildBanRemove           *discordgo.GuildBanRemove
	GuildCreate              *discordgo.GuildCreate
	GuildDelete              *discordgo.GuildDelete
	GuildEmojisUpdate        *discordgo.GuildEmojisUpdate
	GuildMemberAdd           *discordgo.GuildMemberAdd
	GuildMemberRemove        *discordgo.GuildMemberRemove
	GuildMemberUpdate        *discordgo.GuildMemberUpdate
	GuildMembersChunk        *discordgo.GuildMembersChunk
	GuildRoleCreate          *discordgo.GuildRoleCreate
	GuildRoleDelete          *discordgo.GuildRoleDelete
	GuildRoleUpdate          *discordgo.GuildRoleUpdate
	GuildUpdate              *discordgo.GuildUpdate
	MessageCreate            *discordgo.MessageCreate
	MessageDelete            *discordgo.MessageDelete
	MessageReactionAdd       *discordgo.MessageReactionAdd
	MessageReactionRemove    *discordgo.MessageReactionRemove
	MessageReactionRemoveAll *discordgo.MessageReactionRemoveAll
	MessageUpdate            *discordgo.MessageUpdate
	PresenceUpdate           *discordgo.PresenceUpdate
}

// DestinationType is a type for destination types
type DestinationType string

// defines possible destination types
const (
	LambdaDestinationType DestinationType = "lambda"
	SqsDestinationType                    = "sqs"
)

// ErrorHandlerType is the Error Handler Type used for the EventContainer
type ErrorHandlerType string

// defines ErrorHandler types
const (
	SentryErrorHandler  ErrorHandlerType = "sentry"
	DiscordErrorHandler                  = "discord"
)

// DestinationData contains all information for one destination
type DestinationData struct {
	Type          DestinationType
	Name          string
	ErrorHandlers []ErrorHandlerType
	Alias         string
}
