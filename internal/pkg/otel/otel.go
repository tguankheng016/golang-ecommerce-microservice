package otel

import (
	"context"
	"fmt"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/environment"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewJaegerExporter(ctx context.Context, cfg *JaegerOptions) (*otlptrace.Exporter, error) {
	headers := map[string]string{
		"content-type": "application/json",
	}

	exporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, err
	}

	return exporter, nil
}

func NewTraceProvider(cfg *JaegerOptions, exporter *otlptrace.Exporter, env environment.Environment) (trace.Tracer, *tracesdk.TracerProvider, error) {
	tracerProvider := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(cfg.ServiceName),
				attribute.String("environment", env.GetEnvironmentName()),
			),
		),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))

	tracer := tracerProvider.Tracer(cfg.TracerName)

	return tracer, tracerProvider, nil
}

func RunTracing(lc fx.Lifecycle, logger *zap.Logger, tracerProvider *tracesdk.TracerProvider) error {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("closing tracer provider...")

			if err := tracerProvider.Shutdown(ctx); err != nil {
				logger.Info("error when closing tracer provider...", zap.Error(err))
			}

			logger.Info("tracer provider closed gracefully")

			return nil
		},
	})

	return nil
}
