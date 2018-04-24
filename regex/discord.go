package regex

import "regexp"

var (
	// MentionRegex matches Discord User Mentions
	// Source: https://github.com/b1naryth1ef/disco/blob/master/disco/bot/command.py#L15, modified to accept IDs
	MentionRegex = regexp.MustCompile("(<@!?)?([0-9]+)(>)?")

	// RoleRegex matches Discord Role Mentions
	// Source: https://github.com/b1naryth1ef/disco/blob/master/disco/bot/command.py#L16, modified to accept IDs
	RoleRegex = regexp.MustCompile("(<@&)?([0-9]+)(>)?")

	// ChannelRegex matches Discord Channel Mentions
	// Source: https://github.com/b1naryth1ef/disco/blob/master/disco/bot/command.py#L17, modified to accept IDs
	ChannelRegex = regexp.MustCompile("(<#)?([0-9]+)(>)?")
)
