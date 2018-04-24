package mdb

import (
	"strings"
)

// ErrNotFound returns true if the given error is a not found error from MongoDB
// includes errors from invalid object IDs
func ErrNotFound(err error) (notFound bool) {
	if err != nil {
		if strings.Contains(err.Error(), "not found") ||
			strings.Contains(err.Error(), "ObjectIDs must be exactly 12 bytes long") {
			return true
		}
	}
	return false
}
