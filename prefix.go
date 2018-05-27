package dhelpers

var (
	defaultPrefix = "/"
)

// PrefixRule stores the Prefix Config for a Guild
type PrefixRule struct {
	GuildID string
	Prefix  []string
}

// GetPrefixes returns all customized prefix entries for guilds
func GetPrefixes() (prefixRules []PrefixRule) {
	// TODO: load config from MongoDB
	return []PrefixRule{ // example prefix config
		{
			GuildID: "435420687906111498", // Cacophony cord
			Prefix:  []string{"!"},
		},
	}
}

// GetAllPrefix returns all possible prefixes for a specific guildID
func GetAllPrefix(botUserID, guildID string) (prefixes []string) {
	prefixes = append(prefixes, "<@"+botUserID+">", "<@!"+botUserID+">")

	for _, prefix := range GetPrefixes() {
		if prefix.GuildID == guildID {
			prefixes = append(prefixes, prefix.Prefix...)
			return prefixes
		}
	}

	prefixes = append(prefixes, defaultPrefix)
	return prefixes
}

// GetPrefix returns the default prefix for a specific GuildID
func GetPrefix(guildID string) (prefix string) {
	for _, prefix := range GetPrefixes() {
		if prefix.GuildID == guildID && len(prefix.Prefix) >= 0 {
			return prefix.Prefix[0]
		}
	}

	return defaultPrefix
}
