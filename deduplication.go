package dhelpers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// GetEventKey returns an unique key for a discordgo event for deduplication
func GetEventKey(i interface{}) (key string) {
	switch t := i.(type) {
	case *discordgo.GuildCreate:
		return "cacophony:gateway:event-" + string(GuildCreateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Guild))
	case *discordgo.GuildUpdate:
		return "cacophony:gateway:event-" + string(GuildUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Guild))
	case *discordgo.GuildDelete:
		return "cacophony:gateway:event-" + string(GuildDeleteEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Guild))
	case *discordgo.GuildMemberAdd:
		return "cacophony:gateway:event-" + string(GuildMemberAddEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Member))
	case *discordgo.GuildMemberUpdate:
		return "cacophony:gateway:event-" + string(GuildMemberUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Member))
	case *discordgo.GuildMemberRemove:
		return "cacophony:gateway:event-" + string(GuildMemberRemoveEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Member))
	case *discordgo.GuildMembersChunk:
		return "cacophony:gateway:event-" + string(GuildMembersChunkEventType) + "-" + GetMD5Hash(fmt.Sprintf("%s %v", t.GuildID, t.Members))
	case *discordgo.GuildRoleCreate:
		return "cacophony:gateway:event-" + string(GuildRoleCreateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.GuildRole))
	case *discordgo.GuildRoleUpdate:
		return "cacophony:gateway:event-" + string(GuildRoleUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.GuildRole))
	case *discordgo.GuildRoleDelete:
		return "cacophony:gateway:event-" + string(GuildRoleDeleteEventType) + "-" + GetMD5Hash(fmt.Sprintf("%s %s", t.RoleID, t.GuildID))
	case *discordgo.GuildEmojisUpdate:
		return "cacophony:gateway:event-" + string(GuildEmojisUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%s %v", t.GuildID, t.Emojis))
	case *discordgo.ChannelCreate:
		return "cacophony:gateway:event-" + string(ChannelCreateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Channel))
	case *discordgo.ChannelUpdate:
		return "cacophony:gateway:event-" + string(ChannelUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Channel))
	case *discordgo.ChannelDelete:
		return "cacophony:gateway:event-" + string(ChannelDeleteEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Channel))
	case *discordgo.MessageCreate:
		return "cacophony:gateway:event-" + string(MessageCreateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Message))
	case *discordgo.MessageUpdate:
		return "cacophony:gateway:event-" + string(MessageUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Message))
	case *discordgo.MessageDelete:
		return "cacophony:gateway:event-" + string(MessageDeleteEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.Message))
	case *discordgo.PresenceUpdate:
		return "cacophony:gateway:event-" + string(PresenceUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v %s %v", t.Presence, t.GuildID, t.Roles))
	case *discordgo.ChannelPinsUpdate:
		return "cacophony:gateway:event-" + string(ChannelPinsUpdateEventType) + "-" + GetMD5Hash(fmt.Sprintf("%s %s", t.LastPinTimestamp, t.ChannelID))
	case *discordgo.GuildBanAdd:
		return "cacophony:gateway:event-" + string(GuildBanAddEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v %s", t.User, t.GuildID))
	case *discordgo.GuildBanRemove:
		return "cacophony:gateway:event-" + string(GuildBanRemoveEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v %s", t.User, t.GuildID))
	case *discordgo.MessageReactionAdd:
		return "cacophony:gateway:event-" + string(MessageReactionAddEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.MessageReaction))
	case *discordgo.MessageReactionRemove:
		return "cacophony:gateway:event-" + string(MessageReactionRemoveEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.MessageReaction))
	case *discordgo.MessageReactionRemoveAll:
		return "cacophony:gateway:event-" + string(MessageReactionRemoveAllEventType) + "-" + GetMD5Hash(fmt.Sprintf("%v", t.MessageReaction))
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
