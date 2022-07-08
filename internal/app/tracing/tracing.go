package tracing

import (
	"context"
	"github.com/kevinmichaelchen/api-go-template/internal/app/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.uber.org/fx"
	"log"
	"time"
)

const (
	serviceName = "foo-api"
	environment = "production"
	id          = 1
)

var Module = fx.Module("tracing",
	fx.Provide(
		NewTracerProvider,
	),
	fx.Invoke(Register),
)

func Register(tp *tracesdk.TracerProvider) {
	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
}

// NewTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func NewTracerProvider(lc fx.Lifecycle, cfg *config.Config) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.TraceConfig.URL)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always sample traces.
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Do not make the application hang when it is shutdown.
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			log.Println("Shutting down TracerProvider...")
			err := tp.Shutdown(ctx)
			if err != nil {
				return err
			}
			log.Println("Successfully shut down TracerProvider")
			return nil
		},
	})
	return tp, nil
}
