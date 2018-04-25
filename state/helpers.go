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