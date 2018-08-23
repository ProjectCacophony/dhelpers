package mongo

import "errors"

// ErrNotFound is used when the object(s) was not found
var ErrNotFound = errors.New("object(s) not found")

// ErrUnavailable is used when the database or collection is unavailable
var ErrUnavailable = errors.New("database or collection is unavailable")
