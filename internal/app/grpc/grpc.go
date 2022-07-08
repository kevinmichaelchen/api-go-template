package grpc

import (
	"context"
	"fmt"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/kevinmichaelchen/api-go-template/internal/idl/coop/drivers/foo/v1beta1"
	"github.com/kevinmichaelchen/api-go-template/internal/service"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
)

var Module = fx.Module("grpc",
	fx.Provide(
		NewConnectWrapper,
		NewGRPCServer,
		fx.Annotate(
			NewConnectGoServer,
			// Because the output of NewConnectGoServer returns a *http.ServeMux
			// and because we have other DI functions that return that as well,
			// we have to use ResultTags to disambiguate.
			fx.ResultTags(`name:"connectGoMux"`),
		),
	),
	fx.Invoke(
		RegisterGrpcServer,
		fx.Annotate(
			RegisterConnectGoServer,
			// TODO making this positional seems a lot more brittle than just using a struct with annotated fields
			fx.ParamTags(``, ``, `name:"connectGoMux"`),
		),
	),
)

func RegisterGrpcServer(
	svc *service.Service,
	server *grpc.Server,
) {
	// Register our gRPC server
	v1beta1.RegisterFooServiceServer(server, svc)
	grpc_health_v1.RegisterHealthServer(server, svc)
	reflection.Register(server)
}

func NewGRPCServer(lc fx.Lifecycle, logger *zap.Logger) (*grpc.Server, error) {
	// TODO configure options here
	//var opts grpc.ServerOption
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			getUnaryInterceptors(logger)...,
		),
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger),
		),
	)
	lc.Append(fx.Hook{
		// To mitigate the impact of deadlocks in application startup and
		// shutdown, Fx imposes a time limit on OnStart and OnStop hooks. By
		// default, hooks have a total of 15 seconds to complete. Timeouts are
		// passed via Go's usual context.Context.
		OnStart: func(context.Context) error {
			logger.Info("Starting gRPC server.")
			// TODO make configurable
			address := fmt.Sprintf(":%d", 8080)
			lis, err := net.Listen("tcp", address)
			if err != nil {
				logger.Sugar().Fatal("Failed to listen on address \"%s\"", address, zap.Error(err))
			}
			go s.Serve(lis)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping gRPC server.")
			// GracefulStop stops the gRPC server gracefully. It stops the server from
			// accepting new connections and RPCs and blocks until all the pending RPCs are
			// finished.
			s.GracefulStop()
			return nil
		},
	})
	return s, nil
}
