package state

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/project-d-collab/dhelpers/regex"
)

// UserFromMention finds an user in a mention, can be direct ID input
func UserFromMention(mention string) (*discordgo.User, error) {
	result := regex.MentionRegex.FindStringSubmatch(mention)
	if len(result) == 4 {
		return User(result[2])
	}

	return nil, errors.New("user not found")
}

// ChannelFromMention finds a text channel on the same server in an mention, can be direct ID input
func ChannelFromMention(guildID string, mention string) (*discordgo.Channel, error) {
	result := regex.ChannelRegex.FindStringSubmatch(mention)
	if len(result) == 4 {
		channel, err := Channel(result[2])
		if err != nil {
			return nil, err
		}

		if channel.GuildID != guildID {
			return nil, ErrTargetWrongServer
		}

		if channel.Type != discordgo.ChannelTypeGuildText {
			return nil, ErrTargetWrongType
		}

		return channel, nil
	}

	return nil, ErrStateNotFound
}
