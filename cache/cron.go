package cache

import (
	"sync"

	"github.com/robfig/cron"
)

var (
	cronHandler      *cron.Cron
	cronHandlerMutex sync.RWMutex
)

// SetCron caches a cron handler for future use
func SetCron(s *cron.Cron) {
	cronHandlerMutex.Lock()
	defer cronHandlerMutex.Unlock()

	cronHandler = s
}

// GetCron returns a cached cron handler
func GetCron() *cron.Cron {
	cronHandlerMutex.Lock()
	defer cronHandlerMutex.Unlock()

	return cronHandler
}
