package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	"github.com/kevinmichaelchen/api-go-template/internal/idl/coop/drivers/foo/v1beta1/foov1beta1connect"
	"github.com/kevinmichaelchen/api-go-template/internal/service"
	"github.com/rs/cors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

type registerConnectGoServerInput struct {
	fx.In

	Logger     *zap.Logger
	ConnectSvc *service.ConnectWrapper
	Mux        *http.ServeMux `name:"connectGoMux"`
}

func RegisterConnectGoServer(in registerConnectGoServerInput) {
	// Register our Connect-Go server
	path, handler := foov1beta1connect.NewFooServiceHandler(
		in.ConnectSvc,
		connect.WithInterceptors(getUnaryInterceptorsForConnect(in.Logger)...),
	)
	checker := grpchealth.NewStaticChecker(
		// protoc-gen-connect-go generates package-level constants
		// for these fully-qualified protobuf service names, so we'd be able
		// to reference foov1beta1.FooService as opposed to foo.v1beta1.FooService.
		"coop.drivers.foov1beta1.FooService",
	)
	in.Mux.Handle(grpchealth.NewHandler(checker))
	in.Mux.Handle(path, handler)
}

func NewConnectWrapper(s *service.Service) *service.ConnectWrapper {
	return service.NewConnectWrapper(s)
}

type NewConnectGoServerOutput struct {
	fx.Out

	Mux *http.ServeMux `name:"connectGoMux"`
}

func NewConnectGoServer(lc fx.Lifecycle, logger *zap.Logger, cfg Config) NewConnectGoServerOutput {
	mux := http.NewServeMux()
	address := fmt.Sprintf("%s:%d", cfg.ConnectConfig.Host, cfg.ConnectConfig.Port)
	srv := &http.Server{
		Addr: address,
		// Use h2c so we can serve HTTP/2 without TLS.
		Handler: h2c.NewHandler(
			newCORS().Handler(mux),
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
	return NewConnectGoServerOutput{
		Mux: mux,
	}
}

func newCORS() *cors.Cors {
	// To let web developers play with the demo service from browsers, we need a
	// very permissive CORS setup.
	return cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(origin string) bool {
			// Allow all origins, which effectively disables CORS.
			return true
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			// Content-Type is in the default safelist.
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
		},
	})
}
