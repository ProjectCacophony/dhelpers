package dhelpers

import (
	"crypto/md5" // nolint: gas
	"encoding/hex"
)

// GetMD5Hash returns the md5 hash for a string
func GetMD5Hash(text string) string {
	hasher := md5.New()        // nolint: gas
	hasher.Write([]byte(text)) // nolint: errcheck, gas
	return hex.EncodeToString(hasher.Sum(nil))
}
