package tracing

import (
	"context"
	"fmt"
	"math"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.6.1"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "github.com/k8gb-io/k8gb"
)

type Settings struct {
	Enabled       bool
	Endpoint      string
	SamplingRatio float64
	Commit        string
	AppVersion    string
}

func SetupTracing(ctx context.Context, cfg Settings, log *zerolog.Logger) (func(), trace.Tracer) {
	if !cfg.Enabled {
		log.Info().Msg("OTLP tracing is disabled")
		return func() {}, trace.NewNoopTracerProvider().Tracer(instrumentationName)
	}
	log.Info().Msg("OTLP tracing is ON")
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(cfg.Endpoint),
		otlptracehttp.WithInsecure(),
	)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Err(err).Msg("creating OTLP trace exporter")
	}
	var samplerOption sdktrace.TracerProviderOption
	eps := 0.0001
	if math.Abs(cfg.SamplingRatio-1.0) > eps { // not equal 1.0 (IEEE 754)
		samplerOption = sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SamplingRatio))
		log.Info().Msg(fmt.Sprintf("Tracing: sampling ratio is set to '%.3f'", cfg.SamplingRatio))
	} else {
		log.Info().Msg("Tracing: sampling ratio is not specified, using AlwaysSample")
		samplerOption = sdktrace.WithSampler(sdktrace.AlwaysSample())
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource(cfg.Commit)),
		samplerOption,
	)
	otel.SetTracerProvider(tracerProvider)
	tracer := tracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(cfg.AppVersion),
		trace.WithSchemaURL(semconv.SchemaURL))

	return func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Err(err).Msg("stopping tracer provider")
		}
	}, tracer
}

func newResource(commit string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("k8gb"),
		semconv.ServiceVersionKey.String(commit),
	)
}
