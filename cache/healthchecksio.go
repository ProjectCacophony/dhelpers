package cache

import (
	"sync"

	"github.com/kristofferahl/go-healthchecksio"
)

var (
	healthchecksIOClient      *healthchecksio.Client
	healthchecksIOClientMutex sync.RWMutex
)

// SetHealthchecksIO caches an healthchecks.io client for future use
func SetHealthchecksIO(c *healthchecksio.Client) {
	healthchecksIOClientMutex.Lock()
	defer healthchecksIOClientMutex.Unlock()

	healthchecksIOClient = c
}

// GetHealthchecksIO returns a cached healthchecks.io client
func GetHealthchecksIO() *healthchecksio.Client {
	healthchecksIOClientMutex.RLock()
	defer healthchecksIOClientMutex.RUnlock()

	return healthchecksIOClient
}
