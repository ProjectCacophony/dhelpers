package components

import (
	"io"

	ddtrace "github.com/DataDog/dd-trace-go/opentracing"
	"github.com/opentracing/opentracing-go"
	"gitlab.com/Cacophony/dhelpers/cache"
)

var tracerCloser io.Closer

// InitTracer initialises the opentracing tracer and datadog integration
func InitTracer(service string) error {
	var err error
	var tracer opentracing.Tracer

	config := ddtrace.NewConfiguration()
	config.ServiceName = service

	tracer, tracerCloser, err = ddtrace.NewTracer(config)
	if err != nil {
		return err
	}

	// set the Datadog tracer as a cached tracer
	cache.SetTracer(tracer)
	// set the Datadog tracer as a GlobalTracer
	opentracing.SetGlobalTracer(tracer)
	return nil
}

// UninitTracer uninitialises the opentracing tracer and datadog integration
func UninitTracer() error {
	if tracerCloser != nil {
		return tracerCloser.Close()
	}
	return nil
}
