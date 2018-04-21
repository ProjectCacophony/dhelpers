package state

import (
	"github.com/bwmarrin/discordgo"
	"github.com/json-iterator/go"
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
	return readStateSet(guildIdsSetKey())
}

// AllUserIDs returns a list of all User IDs from the shared state
func AllUserIDs() (userIDs []string, err error) {
	return readStateSet(userIdsSetKey())
}
