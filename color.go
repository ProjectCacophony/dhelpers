package dhelpers

import (
	"math/big"
	"strings"
)

// GetDiscordColorFromHex returns discord color integer from a hex code
func GetDiscordColorFromHex(hex string) int {
	colorInt, ok := new(big.Int).SetString(strings.Replace(hex, "#", "", 1), 16)
	if ok {
		return int(colorInt.Int64())
	}

	return 0x0FADED
}

// GetHexFromDiscordColor returns hex code from a discord color integer
func GetHexFromDiscordColor(colour int) (hex string) {
	return strings.ToUpper(big.NewInt(int64(colour)).Text(16))
}
