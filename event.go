package dhelpers

import (
	"context"
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

// ContextKey is a type for keys to use in context.Context
type ContextKey string

const (
	// ContextKeyEvent is the key for cacophony events used in context.Context
	ContextKeyEvent ContextKey = "event"
)

// ContextWithEvent creates a context with a cacophony event item
func ContextWithEvent(event EventContainer) context.Context {
	return context.WithValue(context.Background(), ContextKeyEvent, event)
}

// EventFromContext reads a cacophony event from a context, returns nil if unable to get event
func EventFromContext(ctx context.Context) *EventContainer {
	if ctx == nil {
		return nil
	}

	event, _ := ctx.Value(ContextKeyEvent).(EventContainer)
	return &event
}

// DestinationType is a type for destination types
type DestinationType string

// defines possible destination types
const (
	LambdaDestinationType DestinationType = "lambda"
	SqsDestinationType                    = "sqs"
)

// DestinationData contains all information for one destination
type DestinationData struct {
	Type          DestinationType
	Name          string
	ErrorHandlers []ErrorHandlerType
	Alias         string
}

// CreateEventContainer creates an EventContainer from a discord event
func CreateEventContainer(gatewayStartedAt, receivedAt time.Time, session *discordgo.Session, eventKey string, i interface{}) (container EventContainer) {
	// create enhanced Event
	dDEvent := EventContainer{
		GatewayStarted: gatewayStartedAt,
		ReceivedAt:     receivedAt,
		Key:            eventKey,
		BotUserID:      session.State.User.ID,
	}

	switch t := i.(type) {
	case *discordgo.GuildCreate:
		dDEvent.Type = GuildCreateEventType
		dDEvent.GuildCreate = t
	case *discordgo.GuildUpdate:
		dDEvent.Type = GuildUpdateEventType
		dDEvent.GuildUpdate = t
	case *discordgo.GuildDelete:
		dDEvent.Type = GuildDeleteEventType
		dDEvent.GuildDelete = t
	case *discordgo.GuildMemberAdd:
		dDEvent.Type = GuildMemberAddEventType
		dDEvent.GuildMemberAdd = t
	case *discordgo.GuildMemberUpdate:
		dDEvent.Type = GuildMemberUpdateEventType
		dDEvent.GuildMemberUpdate = t
	case *discordgo.GuildMemberRemove:
		dDEvent.Type = GuildMemberRemoveEventType
		dDEvent.GuildMemberRemove = t
	case *discordgo.GuildMembersChunk:
		dDEvent.Type = GuildMembersChunkEventType
		dDEvent.GuildMembersChunk = t
	case *discordgo.GuildRoleCreate:
		dDEvent.Type = GuildRoleCreateEventType
		dDEvent.GuildRoleCreate = t
	case *discordgo.GuildRoleUpdate:
		dDEvent.Type = GuildRoleUpdateEventType
		dDEvent.GuildRoleUpdate = t
	case *discordgo.GuildRoleDelete:
		dDEvent.Type = GuildRoleDeleteEventType
		dDEvent.GuildRoleDelete = t
	case *discordgo.GuildEmojisUpdate:
		dDEvent.Type = GuildEmojisUpdateEventType
		dDEvent.GuildEmojisUpdate = t
	case *discordgo.ChannelCreate:
		dDEvent.Type = ChannelCreateEventType
		dDEvent.ChannelCreate = t
	case *discordgo.ChannelUpdate:
		dDEvent.Type = ChannelUpdateEventType
		dDEvent.ChannelUpdate = t
	case *discordgo.ChannelDelete:
		dDEvent.Type = ChannelDeleteEventType
		dDEvent.ChannelDelete = t
	case *discordgo.MessageCreate:
		dDEvent.Type = MessageCreateEventType
		// args and prefix
		prefixes := GetAllPrefix(dDEvent.BotUserID, t.GuildID)
		args, prefix := GetMessageArguments(t.Content, prefixes)
		dDEvent.Args = args
		dDEvent.Prefix = prefix
		dDEvent.MessageCreate = t
	case *discordgo.MessageUpdate:
		dDEvent.Type = MessageUpdateEventType
		// args and prefix
		prefixes := GetAllPrefix(dDEvent.BotUserID, t.GuildID)
		args, prefix := GetMessageArguments(t.Content, prefixes)
		dDEvent.Args = args
		dDEvent.Prefix = prefix
		dDEvent.MessageUpdate = t
	case *discordgo.MessageDelete:
		dDEvent.Type = MessageDeleteEventType
		dDEvent.MessageDelete = t
	case *discordgo.ChannelPinsUpdate:
		dDEvent.ChannelPinsUpdate = t
	case *discordgo.GuildBanAdd:
		dDEvent.GuildBanAdd = t
	case *discordgo.GuildBanRemove:
		dDEvent.GuildBanRemove = t
	case *discordgo.MessageReactionAdd:
		dDEvent.MessageReactionAdd = t
	case *discordgo.MessageReactionRemove:
		dDEvent.MessageReactionRemove = t
	case *discordgo.MessageReactionRemoveAll:
		dDEvent.MessageReactionRemoveAll = t
	}

	return dDEvent
}
