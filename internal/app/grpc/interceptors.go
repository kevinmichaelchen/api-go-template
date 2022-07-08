package grpc

import (
	"context"
	"github.com/bufbuild/connect-go"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/kevinmichaelchen/api-go-template/pkg/grpc/interceptors/stats"
	"github.com/kevinmichaelchen/api-go-template/pkg/grpc/interceptors/tracelog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func getUnaryInterceptors(logger *zap.Logger) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		// TODO is it possible not to sample Health/Check calls?
		// Starts a new span
		otelgrpc.UnaryServerInterceptor(),
		// Adds logger to context
		grpc_zap.UnaryServerInterceptor(logger),
		// Add trace ID as field on logger
		tracelog.UnaryServerInterceptor(),
		// Response counts (w/ status code as a dimension)
		stats.UnaryServerInterceptor(),
	}
}

func getUnaryInterceptorsForConnect(logger *zap.Logger) []connect.Interceptor {
	return []connect.Interceptor{
		connectInterceptorForSpan(),
		connectInterceptorForLogger(logger),
		// Add trace ID as field on logger
		tracelog.UnaryServerInterceptorForConnect(),
		// Response counts (w/ status code as a dimension)
		stats.UnaryServerInterceptorForConnect(),
	}
}

func connectInterceptorForSpan() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			name := req.Spec().Procedure
			tr := otel.Tracer("")
			ctx, span := tr.Start(ctx, name)
			defer span.End()

			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func connectInterceptorForLogger(logger *zap.Logger) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			return next(ctxzap.ToContext(ctx, logger), req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
