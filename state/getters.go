package state

import (
	"github.com/bwmarrin/discordgo"
	jsoniter "github.com/json-iterator/go"
)

func Guild(guildID string) (guild *discordgo.Guild, err error) {
	data, err := readStateObject(guildKey(guildID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &guild)
	return
}

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

func Member(guildID, userID string) (member *discordgo.Member, err error) {
	data, err := readStateObject(memberKey(guildID, userID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &member)
	return
}

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

func Channel(channelID string) (channel *discordgo.Channel, err error) {
	data, err := readStateObject(channelKey(channelID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &channel)
	return
}

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

func User(userID string) (user *discordgo.User, err error) {
	data, err := readStateObject(userKey(userID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &user)
	return
}
