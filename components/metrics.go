package components

import (
	"time"

	"gitlab.com/project-d-collab/dhelpers/metrics"
)

// InitMetrics initializes metrics and sets Uptime
func InitMetrics() {
	metrics.Uptime.Set(time.Now().Unix())
}
