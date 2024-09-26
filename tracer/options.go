package tracer

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type OptionTracer func(t *optTracer)

type optTracer struct {
	// Service name
	ServiceName string

	// Agent Tracer Host Platform
	AgentHost string

	// Service version
	VersionService string

	// Service instance id
	ServiceInstanceId string

	// Agent Tracer Platform
	AgentPlatform PlatformType

	// RatioSample will record the ratio of showing traces
	RatioSampler float64

	// Collector tls to secure communications recomended if using Open-Telemetry Collector
	TLS *tls.Config
}

// FinishOption for an option when trace is finished
type FinishOption struct {
	Tags  map[string]interface{}
	Error error
}

// FinishOptionFunc func
type FinishOptionFunc func(*FinishOption)

func defaultOptionTracer() *optTracer {
	return &optTracer{
		ServiceName:       "default-service-name",
		AgentHost:         "http://localhost:4317",
		AgentPlatform:     OTLP,
		VersionService:    "1.0.0",
		ServiceInstanceId: uuid.New().String(),
		RatioSampler:      5.0,
		TLS:               nil,
	}
}

// Option Tracer

func WithTracerAgentPlatform(exporterPlatform PlatformType) OptionTracer {
	return func(t *optTracer) {
		t.AgentPlatform = exporterPlatform
	}
}

func WithTracerServiceName(name string) OptionTracer {
	return func(t *optTracer) {
		t.ServiceName = name
	}
}

func WithTracerAgentHost(uri string) OptionTracer {
	return func(t *optTracer) {
		t.AgentHost = uri
	}
}

func WithTracerTLS(tls *tls.Config) OptionTracer {
	return func(t *optTracer) {
		t.TLS = tls
	}
}

func WithServiceVersion(version string) OptionTracer {
	return func(t *optTracer) {
		t.VersionService = version
	}
}

func WithServiceInstanceId(instanceId string) OptionTracer {
	return func(t *optTracer) {
		t.ServiceInstanceId = instanceId
	}
}

func WithRatioSampler(ratio float64) OptionTracer {
	return func(t *optTracer) {
		t.RatioSampler = ratio
	}
}

func toValue(v interface{}) string {
	var str string
	switch val := v.(type) {

	case uint, uint64, int, int64, float32, float64:
		str = fmt.Sprintf("%v", val)
	case error:
		if val != nil {
			str = val.Error()
		}
	case string:
		str = val
	case []byte:
		str = string(val)
	default:
		b, _ := json.Marshal(val)
		str = string(b)
	}

	return str
}
