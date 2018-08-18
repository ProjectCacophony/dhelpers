package components

import (
	"os"

	"github.com/kristofferahl/go-healthchecksio"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// InitHealthchecksIO initializes and caches the healthchecks.io client
// reads the healthchecks.io api key from the environment variable Healthchecks.IO
func InitHealthchecksIO() {
	// create a new Healthchecks.IO API Client
	healthchecksIOClient := healthchecksio.NewClient(
		os.Getenv("HEALTHCHECKSIO_API_KEY"),
	)

	// cache client
	cache.SetHealthchecksIO(healthchecksIOClient)
}
