package state

import (
	"github.com/bwmarrin/discordgo"
	"github.com/json-iterator/go"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// Guild returns the specified Guild from the shard state, returns ErrStateNotFound if not found
func Guild(guildID string) (guild *discordgo.Guild, err error) {
	data, err := readStateObject(guildKey(guildID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &guild)
	return
}

// Presence returns the specified Presence from the shard state, returns ErrStateNotFound if not found
func Presence(guildID, userID string) (presence *discordgo.Presence, err error) {
	guild, err := Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, presence := range guild.Presences {
		if presence.User.ID == userID {
			return presence, nil
		}
	}

	return nil, ErrStateNotFound
}

// Member returns the specified Member from the shard state, returns ErrStateNotFound if not found
func Member(guildID, userID string) (member *discordgo.Member, err error) {
	data, err := readStateObject(memberKey(guildID, userID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &member)
	return
}

// Role returns the specified Role from the shard state, returns ErrStateNotFound if not found
func Role(guildID, roleID string) (role *discordgo.Role, err error) {
	guild, err := Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, role := range guild.Roles {
		if role.ID == roleID {
			return role, nil
		}
	}

	return nil, ErrStateNotFound
}

// Channel returns the specified Channel from the shard state, returns ErrStateNotFound if not found
func Channel(channelID string) (channel *discordgo.Channel, err error) {
	data, err := readStateObject(channelKey(channelID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &channel)
	return
}

// Emoji returns the specified Emoji from the shard state, returns ErrStateNotFound if not found
func Emoji(guildID, emojiID string) (emoji *discordgo.Emoji, err error) {
	guild, err := Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, emoji := range guild.Emojis {
		if emoji.ID == emojiID {
			return emoji, nil
		}
	}

	return nil, ErrStateNotFound
}

// User returns the specified User from the shard state, returns ErrStateNotFound if not found
func User(userID string) (user *discordgo.User, err error) {
	data, err := readStateObject(userKey(userID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &user)
	return
}

// AllGuildIDs returns a list of all Guild IDs from the shared state
func AllGuildIDs() (guildIDs []string, err error) {
	return readStateSet(allGuildIdsSetKey())
}

// AllUserIDs returns a list of all User IDs from the shared state
func AllUserIDs() (userIDs []string, err error) {
	return readStateSet(allUserIdsSetKey())
}

// GuildUserIDs returns a list of all User IDs in a specific Guild from the shared state
func GuildUserIDs(guildID string) (userIDs []string, err error) {
	return readStateSet(guildUserIdsSetKey(guildUserIdsSetKey(guildID)))
}

// IsMember true if the User is a member of the specified Guild
func IsMember(guildID, userID string) (isMember bool, err error) {
	isMember, err = cache.GetRedisClient().SIsMember(guildUserIdsSetKey(guildID), userID).Result()
	return isMember, err
}

// UserChannelPermissions returns the permission of a user in a channel
func UserChannelPermissions(userID, channelID string) (apermissions int, err error) {
	var channel *discordgo.Channel
	channel, err = Channel(channelID)
	if err != nil {
		return
	}

	var guild *discordgo.Guild
	guild, err = Guild(channel.GuildID)
	if err != nil {
		return
	}

	if userID == guild.OwnerID {
		apermissions = discordgo.PermissionAll
		return
	}

	var member *discordgo.Member
	member, err = Member(guild.ID, userID)
	if err != nil {
		return
	}

	return memberPermissions(guild, channel, member), nil
}

// MemberPermissions calculates the permissions for a member
// Source: https://github.com/bwmarrin/discordgo/blob/develop/restapi.go#L503
func memberPermissions(guild *discordgo.Guild, channel *discordgo.Channel, member *discordgo.Member) (apermissions int) {
	userID := member.User.ID

	if userID == guild.OwnerID {
		apermissions = discordgo.PermissionAll
		return
	}

	for _, role := range guild.Roles {
		if role.ID == guild.ID {
			apermissions |= role.Permissions
			break
		}
	}

	for _, role := range guild.Roles {
		for _, roleID := range member.Roles {
			if role.ID == roleID {
				apermissions |= role.Permissions
				break
			}
		}
	}

	if apermissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		apermissions |= discordgo.PermissionAll
	}

	// Apply @everyone overrides from the channel.
	for _, overwrite := range channel.PermissionOverwrites {
		if guild.ID == overwrite.ID {
			apermissions &= ^overwrite.Deny
			apermissions |= overwrite.Allow
			break
		}
	}

	denies := 0
	allows := 0

	// Member overwrites can override role overrides, so do two passes
	for _, overwrite := range channel.PermissionOverwrites {
		for _, roleID := range member.Roles {
			if overwrite.Type == "role" && roleID == overwrite.ID {
				denies |= overwrite.Deny
				allows |= overwrite.Allow
				break
			}
		}
	}

	apermissions &= ^denies
	apermissions |= allows

	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.Type == "member" && overwrite.ID == userID {
			apermissions &= ^overwrite.Deny
			apermissions |= overwrite.Allow
			break
		}
	}

	if apermissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		apermissions |= discordgo.PermissionAllChannel
	}

	return apermissions
}
