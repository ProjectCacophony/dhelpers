package dhelpers

import (
	"math/big"
	"strings"
)

// HexToDecimal returns a decimal color from a hex color
func HexToDecimal(hex string) int {
	colorInt, ok := new(big.Int).SetString(strings.Replace(hex, "#", "", 1), 16)
	if ok {
		return int(colorInt.Int64())
	}

	return 0x0FADED
}

// DecimalToHex returns a hex color form a decimal color
func DecimalToHex(colour int) (hex string) {
	if colour < 0 {
		return "FADED"
	}

	return strings.ToUpper(big.NewInt(int64(colour)).Text(16))
}
