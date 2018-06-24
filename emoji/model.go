package emoji

import "gitlab.com/Cacophony/dhelpers/random"

// Emoji is a list of emoji strings with helper methods
type Emoji []string

// Get returns a random emoji for the type with < > around it
// Example: <:robyulblush:327206930437373952>
func (e Emoji) Get() string {
	// skip random choice if just one
	if len(e) == 1 {
		return "<" + e[0] + ">"
	}

	// get random emoji
	return "<" + random.FromStringSlice(e) + ">"
}

// GetWithout returns a random emoji for the type without < > around it
// Example: :robyulblush:327206930437373952
func (e Emoji) GetWithout() string {
	// skip random choice if just one
	if len(e) == 1 {
		return e[0]
	}

	// get random emoji
	return random.FromStringSlice(e)
}
