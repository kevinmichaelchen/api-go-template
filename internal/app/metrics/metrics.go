package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
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

	Exporter *otelprom.Exporter
	Mux      *http.ServeMux `name:"metricsMux"`
}

func Register(in registerInput) error {
	// Set global meter provider
	provider := metric.NewMeterProvider(metric.WithReader(in.Exporter))
	global.SetMeterProvider(provider)

	// Register the Prometheus export handler on our Mux HTTP Server.
	registry := prometheus.NewRegistry()
	err := registry.Register(in.Exporter.Collector)
	if err != nil {
		return err
	}
	var handler http.Handler = promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	in.Mux.HandleFunc("/", handler.ServeHTTP)
	return nil
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

func NewPrometheusExporter(lc fx.Lifecycle) (*otelprom.Exporter, error) {
	exporter := otelprom.New()
	return &exporter, nil
}
