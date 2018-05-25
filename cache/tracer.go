package cache

import (
	"sync"

	"github.com/opentracing/opentracing-go"
)

var (
	tracer      opentracing.Tracer
	tracerMutex sync.RWMutex
)

// SetTracer caches a tracer for future use
func SetTracer(s opentracing.Tracer) {
	tracerMutex.Lock()
	tracer = s
	tracerMutex.Unlock()
}

// GetTracer returns a cached tracer
func GetTracer() opentracing.Tracer {
	tracerMutex.RLock()
	defer tracerMutex.RUnlock()

	return tracer
}
