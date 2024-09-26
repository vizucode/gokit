package tracer

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type otplTracePlatform struct {

	// Tracer
	ctx  context.Context
	span trace.Span
	tags map[string]interface{}
}

func initConn(opt *optTracer) *grpc.ClientConn {
	var creds credentials.TransportCredentials

	switch opt.TLS {
	case nil:
		creds = insecure.NewCredentials()
	default:
		creds = credentials.NewTLS(opt.TLS)
	}

	conn, err := grpc.NewClient(
		opt.AgentHost,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatal("failed to create gRPC connection to collector: %w", err)
	}

	return conn
}

func initOTLP(opts *optTracer) (*sdktrace.TracerProvider, Platform) {
	ctx := context.Background()

	// Set ini connection gRPC
	conn := initConn(opts)

	// Create Exporter with Open-Telemetry Collector
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatal("failed to create resource: %w", err)
	}

	// Register the Open-Telemetry Exporter with a TracerProvider
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)

	// set tracer provider
	ot := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(opts.RatioSampler)),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(opts.ServiceName),
				semconv.ServiceVersion(opts.VersionService),
			),
		),
		sdktrace.WithSpanProcessor(bsp),
	)

	return ot, &otplTracePlatform{}
}

func (ot *otplTracePlatform) Start(ctx context.Context, operationName string) Tracer {
	var span trace.Span
	ctx, span = otel.Tracer("otlp").Start(ctx, operationName)
	if span != nil {
		span = trace.SpanFromContext(ctx)
		ctx = trace.ContextWithSpan(ctx, span)
	}

	return &otplTracePlatform{span: span, ctx: ctx}
}

func (otlp *otplTracePlatform) GetTraceID(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).TraceID().String()
}

func (otlp *otplTracePlatform) GetSpanID(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).SpanID().String()
}

func (otlp *otplTracePlatform) Context() (ctx context.Context) {
	return otlp.ctx
}

func (otlp *otplTracePlatform) NewContext() (ctx context.Context) {
	return trace.ContextWithSpan(context.Background(), otlp.span)
}

func (otlp *otplTracePlatform) Tags() (resp map[string]interface{}) {
	return otlp.tags
}

func (otlp *otplTracePlatform) SetTag(key string, value interface{}) {
	if otlp.span == nil {
		return
	}

	if otlp.tags == nil {
		otlp.tags = make(map[string]interface{})
	}

	otlp.tags[key] = value
}

func (otlp *otplTracePlatform) Log(key string, value interface{}) {
	if otlp.span == nil {
		return
	}

	otlp.span.AddEvent("", trace.WithAttributes(attribute.String(key, toValue(value))))
}

func (otlp *otplTracePlatform) SetError(err error) {
	if otlp.span == nil || err == nil {
		return
	}

	otlp.span.RecordError(err)
	otlp.span.SetAttributes(attribute.String("error.message", err.Error()))
}

func (otlp *otplTracePlatform) Finish(opts ...FinishOptionFunc) {
	if otlp.span == nil {
		return
	}

	var finishOpt FinishOption
	for _, opt := range opts {
		if opt != nil {
			opt(&finishOpt)
		}
	}

	if finishOpt.Tags != nil && otlp.tags == nil {
		otlp.tags = make(map[string]interface{})
	}

	for k, v := range finishOpt.Tags {
		otlp.tags[k] = v
	}

	for k, v := range otlp.tags {
		otlp.span.SetAttributes(attribute.String(k, toValue(v)))
	}

	if finishOpt.Error != nil {
		otlp.span.RecordError(finishOpt.Error)
		otlp.span.SetAttributes(attribute.String("error.message", finishOpt.Error.Error()))
	}

	otlp.span.End()
}
