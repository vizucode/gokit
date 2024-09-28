package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

/*
	Tracer kit, used send data to exporter
*/

type Tracer interface {
	Context() context.Context
	NewContext() context.Context
	Tags() map[string]interface{}
	SetTag(key string, value interface{})
	Log(key string, value interface{})
	SetError(err error)
	Finish(opts ...FinishOptionFunc)
}

func New(ServiceName string, opts ...OptionTracer) {
	var (
		platform Platform
		tracer   *sdktrace.TracerProvider
	)

	tracerObj := defaultOptionTracer()

	for _, o := range opts {
		o(tracerObj)
	}

	switch tracerObj.AgentPlatform {
	case OTLP:
		tracer, platform = initOTLP(tracerObj)
	}

	// Set global Tracer Provider
	otel.SetTracerProvider(tracer)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})
	SetTracerPlatformType(platform)
}

// StartTrace starting trace child span from parent span
func StartTrace(ctx context.Context, operationName string) Tracer {
	if activeTracer == nil {
		return nil
	}

	return activeTracer.Start(ctx, operationName)
}

// StartTraceWithContext starting trace child span from parent span with context
func StartTraceWithContext(ctx context.Context, operationName string) (Tracer, context.Context) {
	t := StartTrace(ctx, operationName)

	if t == nil {
		return nil, ctx
	}

	return t, t.Context()
}

// GetTraceID get active trace id
func GetTraceID(ctx context.Context) string {
	if activeTracer == nil {
		return ""
	}

	return activeTracer.GetTraceID(ctx)
}

// GetSpanID get current span id
func GetSpanID(ctx context.Context) string {
	if activeTracer == nil {
		return ""
	}

	return activeTracer.GetSpanID(ctx)
}
