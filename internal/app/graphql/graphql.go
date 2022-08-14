package graphql

import (
	"context"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/rs/cors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

var Module = fx.Module("graphql",
	fx.Provide(
		NewConfig,
		NewGraphQL,
		NewSchema,
	),
	fx.Invoke(
		RegisterGraphQL,
	),
)

type registerGraphQLInput struct {
	fx.In

	Schema *graphql.Schema
	Logger *zap.Logger
	Mux    *http.ServeMux `name:"graphqlMux"`
}

type NewGraphQLOutput struct {
	fx.Out

	Mux *http.ServeMux `name:"graphqlMux"`
}

func NewGraphQL(lc fx.Lifecycle, logger *zap.Logger, cfg Config, crs *cors.Cors) NewGraphQLOutput {
	mux := http.NewServeMux()
	address := fmt.Sprintf("%s:%d", cfg.GraphQLConfig.Host, cfg.GraphQLConfig.Port)
	srv := &http.Server{
		Addr: address,
		// Use h2c so we can serve HTTP/2 without TLS.
		Handler: h2c.NewHandler(
			crs.Handler(mux),
			&http2.Server{},
		),
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go func() {
				err := srv.ListenAndServe()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("connect-go ListenAndServe failed", zap.Error(err))
				}
			}()
			logger.Sugar().Infof("Listening for connect-go on: %s", address)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return NewGraphQLOutput{
		Mux: mux,
	}
}

func RegisterGraphQL(in registerGraphQLInput) {
	h := handler.New(&handler.Config{
		Schema:   in.Schema,
		Pretty:   true,
		GraphiQL: true,
	})
	in.Mux.Handle("/graphql", h)
}
