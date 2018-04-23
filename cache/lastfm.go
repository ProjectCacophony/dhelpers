package cache

import (
	"sync"

	"github.com/Seklfreak/lastfm-go/lastfm"
)

var (
	lastfmClient      *lastfm.Api
	lastfmClientMutex sync.RWMutex
)

// SetLastfFm caches a lastfm client for future use
func SetLastfFm(s *lastfm.Api) {
	lastfmClientMutex.Lock()
	defer lastfmClientMutex.Unlock()

	lastfmClient = s
}

// GetLastFm returns a cached lastfm client
func GetLastFm() *lastfm.Api {
	lastfmClientMutex.Lock()
	defer lastfmClientMutex.Unlock()

	return lastfmClient
}
