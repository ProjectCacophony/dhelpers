package dhelpers

// PrefixRule stores the Prefix Config for a Guild
type PrefixRule struct {
	GuildID string
	Prefix  []string
}

// GetPrefixes returns all customized prefix entries for guilds
func GetPrefixes() (prefixRules []PrefixRule, err error) {
	// TODO: load config from MongoDB
	return []PrefixRule{ // example prefix config
		{
			GuildID: "435420687906111498", // Project D cord
			Prefix:  []string{"!"},
		},
	}, nil
}

// GetPrefix returns all possible prefixes for a specific guildID
func GetPrefix(prefixRules []PrefixRule, botUserID, guildID string) (prefixes []string) {
	prefixes = append(prefixes, "<@"+botUserID+">", "<@!"+botUserID+">")

	for _, prefix := range prefixRules {
		if prefix.GuildID == guildID {
			prefixes = append(prefixes, prefix.Prefix...)
			return prefixes
		}
	}

	prefixes = append(prefixes, "/")
	return prefixes
}
