package components

import (
	"github.com/robfig/cron"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// InitCron initializes and caches a cron handler
func InitCron() {
	c := cron.New()
	c.Start()

	cache.SetCron(c)
}
