package metrics

import (
	"context"
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

var Module = fx.Module("metrics",
	fx.Provide(
		NewMux,
		NewPrometheusExporter,
	),
	fx.Invoke(
		Register,
	),
)

type registerInput struct {
	fx.In

	Exporter *prometheus.Exporter
	Mux      *http.ServeMux `name:"metricsMux"`
}

func Register(in registerInput) {
	// Set global meter provider
	global.SetMeterProvider(in.Exporter.MeterProvider())

	// Register the Prometheus export handler on our Mux HTTP Server.
	in.Mux.HandleFunc("/", in.Exporter.ServeHTTP)
}

type NewMuxOutput struct {
	fx.Out

	Mux *http.ServeMux `name:"metricsMux"`
}

func NewMux(lc fx.Lifecycle) NewMuxOutput {
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

	return NewMuxOutput{
		Mux: mux,
	}
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
