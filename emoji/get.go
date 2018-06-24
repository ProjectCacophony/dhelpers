package emoji

import (
	"strings"

	"gitlab.com/Cacophony/dhelpers/cache"
)

// Get returns a random emoji for the type with < > around it, can be with or without colons around it
// Example: robyulblush => <:robyulblush:327206930437373952>
func Get(name string) string {
	emoji, ok := list[strings.ToLower(strings.Trim(name, ":"))]
	if ok {
		// return requested emoji
		return emoji.Get()
	}

	// log warning
	if cache.HasLogger() {
		cache.GetLogger().Warnln("unable to find emoji", name)
	}

	// return input
	return name
}

// GetWithout returns a random emoji for the type without < > around it, can be with or without colons around it
// Example: robyulblush => :robyulblush:327206930437373952
func GetWithout(name string) string {
	emoji, ok := list[strings.ToLower(strings.Trim(name, ":"))]
	if ok {
		// return requested emoji
		return emoji.GetWithout()
	}

	// log warning
	if cache.HasLogger() {
		cache.GetLogger().Warnln("unable to find emoji", name)
	}

	// return input
	return name
}
