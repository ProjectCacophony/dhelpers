package dhelpers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// GetEventKey returns an unique key for a discordgo event for deduplication
func GetEventKey(i interface{}) (key string) {
	switch t := i.(type) {
	case *discordgo.GuildCreate:
		return "project-d:gateway:event-" + string(GuildCreateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Guild))
	case *discordgo.GuildUpdate:
		return "project-d:gateway:event-" + string(GuildUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Guild))
	case *discordgo.GuildDelete:
		return "project-d:gateway:event-" + string(GuildDeleteEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Guild))
	case *discordgo.GuildMemberAdd:
		return "project-d:gateway:event-" + string(GuildMemberAddEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Member))
	case *discordgo.GuildMemberUpdate:
		return "project-d:gateway:event-" + string(GuildMemberUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Member))
	case *discordgo.GuildMemberRemove:
		return "project-d:gateway:event-" + string(GuildMemberRemoveEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Member))
	case *discordgo.GuildMembersChunk:
		return "project-d:gateway:event-" + string(GuildMembersChunkEventType) + "-" + GetMD5Hash(fmt.Sprintf("%s %v", t.GuildID, t.Members))
	case *discordgo.GuildRoleCreate:
		return "project-d:gateway:event-" + string(GuildRoleCreateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.GuildRole))
	case *discordgo.GuildRoleUpdate:
		return "project-d:gateway:event-" + string(GuildRoleUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.GuildRole))
	case *discordgo.GuildRoleDelete:
		return "project-d:gateway:event-" + string(GuildRoleDeleteEventType) + "-" + GetMD5Hash(fmt.Sprintf("%s %s", t.RoleID, t.GuildID))
	case *discordgo.GuildEmojisUpdate:
		return "project-d:gateway:event-" + string(GuildEmojisUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%s %v", t.GuildID, t.Emojis))
	case *discordgo.ChannelCreate:
		return "project-d:gateway:event-" + string(ChannelCreateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Channel))
	case *discordgo.ChannelUpdate:
		return "project-d:gateway:event-" + string(ChannelUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Channel))
	case *discordgo.ChannelDelete:
		return "project-d:gateway:event-" + string(ChannelDeleteEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Channel))
	case *discordgo.MessageCreate:
		return "project-d:gateway:event-" + string(MessageCreateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Message))
	case *discordgo.MessageUpdate:
		return "project-d:gateway:event-" + string(MessageUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Message))
	case *discordgo.MessageDelete:
		return "project-d:gateway:event-" + string(MessageDeleteEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Message))
	case *discordgo.PresenceUpdate:
		return "project-d:gateway:event-" + string(PresenceUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v %s %v", t.Presence, t.GuildID, t.Roles))
	case *discordgo.ChannelPinsUpdate:
		return "project-d:gateway:event-" + string(ChannelPinsUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%s %s", t.LastPinTimestamp, t.ChannelID))
	case *discordgo.GuildBanAdd:
		return "project-d:gateway:event-" + string(GuildBanAddEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v %s", t.User, t.GuildID))
	case *discordgo.GuildBanRemove:
		return "project-d:gateway:event-" + string(GuildBanRemoveEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v %s", t.User, t.GuildID))
	case *discordgo.MessageReactionAdd:
		return "project-d:gateway:event-" + string(MessageReactionAddEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.MessageReaction))
	case *discordgo.MessageReactionRemove:
		return "project-d:gateway:event-" + string(MessageReactionRemoveEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.MessageReaction))
	case *discordgo.MessageReactionRemoveAll:
		return "project-d:gateway:event-" + string(MessageReactionRemoveAllEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.MessageReaction))
	}
	return ""
}

// IsNewEvent returns true if the event key is new, returns false if the event key has already been handled by other gateways
func IsNewEvent(redisClient *redis.Client, source, eventKey string) (new bool) {
	set, err := redisClient.SetNX(eventKey+":"+source, true, time.Minute*5).Result()
	if err != nil {
		cache.GetLogger().Errorln("error doing deduplication:", err.Error())
		return false
	}

	return set
}
