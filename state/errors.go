package state

import (
	"errors"
)

// ErrStateNotFound will be returned if the item was not found in the shared state
var ErrStateNotFound = errors.New("shared state cache not found")
