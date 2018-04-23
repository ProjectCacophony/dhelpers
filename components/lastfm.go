package components

import (
	"os"

	"github.com/Seklfreak/lastfm-go/lastfm"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// InitLastFm initializes and caches the lastfm client
// reads the lastfm api key from the environment variable LASTFM_API_KEY
// reads the lastfm api secret from the environment variable LASTFM_API_SECRET
func InitLastFm() {
	// create a new LastFM API Client
	lastFmClient := lastfm.New(
		os.Getenv("LASTFM_API_KEY"),
		os.Getenv("LASTFM_API_SECRET"),
	)

	// cache client
	cache.SetLastfFm(lastFmClient)
}
