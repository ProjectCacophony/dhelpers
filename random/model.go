package random

import (
	"math/rand"

	"time"
)

type random struct {
	*rand.Rand
}

var randomObject = &random{
	Rand: rand.New(rand.NewSource(time.Now().Unix())),
}

// FromSlice returns a random item from a slice
func FromSlice(slice []interface{}) interface{} {
	return slice[randomObject.Intn(len(slice))]
}

// FromStringSlice returns a random string from a string slice
func FromStringSlice(slice []string) string {

	return slice[randomObject.Intn(len(slice))]
}
