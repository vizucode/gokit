package tracer

import (
	"context"
	"sync"
)

type PlatformType string

const (
	// Open-Telemetry Collector
	OTLP PlatformType = "otlp"
)

var (
	once         sync.Once
	activeTracer Platform
)

// Platform defines the tracing platform
type Platform interface {
	Start(ctx context.Context, operationName string) Tracer
	GetTraceID(ctx context.Context) string
	GetSpanID(ctx context.Context) string
}

// SetTracerPlatformType function for set tracer platform
func SetTracerPlatformType(t Platform) {
	once.Do(
		func() {
			activeTracer = t
		},
	)
}
