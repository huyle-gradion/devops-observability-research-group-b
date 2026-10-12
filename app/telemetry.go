package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

const (
	// serviceName is the OTel resource service.name for this runtime, joining
	// every span in the viewer under one logical service.
	serviceName = "devops-app-huy"

	// defaultOTLPEndpoint is the in-cluster Collector OTLP/HTTP base endpoint
	// used when OTEL_EXPORTER_OTLP_ENDPOINT is unset. Per the OTLP/HTTP spec a
	// base endpoint gets the signal path appended (see otlpTracesPath).
	defaultOTLPEndpoint = "http://otel-collector.monitoring.svc.cluster.local:4318"

	// otlpTracesPath is the OTLP/HTTP traces signal path relative to the base
	// endpoint.
	otlpTracesPath = "/v1/traces"
)

// initTelemetry builds the OpenTelemetry trace pipeline for the resolved
// runtime id, installs the global tracer provider and propagator once, and
// returns a shutdown function that flushes spans and stops the exporter.
//
// The exporter is configured from OTEL_EXPORTER_OTLP_ENDPOINT (a base URL per
// the OTLP spec), defaulting to the in-cluster Collector. Exporting uses OTLP
// HTTP with a batch processor and an always-on sampler suitable for this small
// learning workload. Building the exporter does not contact the Collector: an
// absent or unreachable endpoint surfaces later as an honest asynchronous
// export error, not a fake success. A malformed endpoint is a hard error so the
// caller can treat it as a startup failure.
func initTelemetry(ctx context.Context, id string) (func(context.Context) error, error) {
	exporter, err := newTraceExporter(ctx)
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceInstanceID(id),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("build telemetry resource: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return provider.Shutdown, nil
}

// newTraceExporter builds the OTLP/HTTP trace exporter from the base endpoint.
// WithEndpoint carries host and port, WithURLPath carries the per-signal path
// (the SDK only appends the default path when it is left empty), and a non-HTTPS
// scheme selects plaintext transport. The returned exporter is a
// sdktrace.SpanExporter so the concrete otlptrace type stays an implementation
// detail.
func newTraceExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = defaultOTLPEndpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse OTEL_EXPORTER_OTLP_ENDPOINT %q: %w", endpoint, err)
	}
	if u.Host == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT %q has no host", endpoint)
	}

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(u.Host),
		otlptracehttp.WithURLPath(path.Join(u.Path, otlpTracesPath)),
	}
	if u.Scheme != "https" {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("create OTLP trace exporter: %w", err)
	}
	return exporter, nil
}
