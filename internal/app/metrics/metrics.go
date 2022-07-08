package metrics

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.uber.org/fx"
	"log"
	"net/http"
)

const muxName = "metricsMux"

var Module = fx.Module("metrics",
	fx.Provide(
		fx.Annotate(
			NewMux,
			fx.ResultTags(fmt.Sprintf(`name:"%s"`, muxName)),
		),
		NewPrometheusExporter,
	),
	fx.Invoke(
		fx.Annotate(
			Register,
			fx.ParamTags(``, fmt.Sprintf(`name:"%s"`, muxName)),
		),
	),
)

func Register(exporter *prometheus.Exporter, mux *http.ServeMux) {
	// Set global meter provider
	global.SetMeterProvider(exporter.MeterProvider())

	// Register the Prometheus export handler on our Mux HTTP Server.
	mux.HandleFunc("/", exporter.ServeHTTP)
}

// TODO use https://pkg.go.dev/go.uber.org/fx#hdr-Named_Values
func NewMux(lc fx.Lifecycle) *http.ServeMux {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":2222",
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Println("Starting HTTP Prometheus server...")
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go server.ListenAndServe()
			log.Println("Prometheus server running on :2222")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return mux
}

func NewPrometheusExporter(lc fx.Lifecycle) (*prometheus.Exporter, error) {
	// Configs for the tally reporter.
	config := prometheus.Config{
		DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	return prometheus.New(config, c)
}
