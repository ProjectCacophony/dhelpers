package metrics

import "expvar"

var (
	// Uptime contains the timestamp when the service was started
	Uptime = expvar.NewInt("uptime")
)
