package components

import (
	"github.com/opentracing/opentracing-go"
	"gitlab.com/Cacophony/dhelpers/cache"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// InitTracer initialises the opentracing tracer and datadog integration
func InitTracer(service string) error {
	t := opentracer.New(tracer.WithServiceName(service))

	// set the Datadog tracer as a cached tracer
	cache.SetTracer(t)
	// set the Datadog tracer as a GlobalTracer
	opentracing.SetGlobalTracer(t)
	return nil
}
