package tracelog

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor sets the trace ID on the request logger.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		logger := ctxzap.Extract(ctx)
		span := trace.SpanFromContext(ctx)
		traceID := span.SpanContext().TraceID().String()
		newCtx := ctxzap.ToContext(ctx, logger.With(
			zap.String("traceid", traceID),
		))

		resp, err = handler(newCtx, req)

		return
	}
}

// UnaryServerInterceptorForConnect sets the trace ID on the request logger.
func UnaryServerInterceptorForConnect() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			logger := ctxzap.Extract(ctx)
			span := trace.SpanFromContext(ctx)
			traceID := span.SpanContext().TraceID().String()
			newCtx := ctxzap.ToContext(ctx, logger.With(
				zap.String("traceid", traceID),
			))

			res, err := next(newCtx, req)
			if err != nil {
				return nil, err
			}

			return res, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
