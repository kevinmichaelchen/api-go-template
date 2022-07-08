package stats

import (
	"context"
	"github.com/bufbuild/connect-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		handleStatusMetrics(ctx, err)

		return
	}
}

func UnaryServerInterceptorForConnect() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			res, err := next(ctx, req)

			handleStatusMetrics(ctx, err)

			return res, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func handleStatusMetrics(ctx context.Context, err error) {
	meter := global.Meter("go.opentelemetry.io/otel/exporters/prometheus")
	counter, err := meter.SyncFloat64().Counter("ex.com.three")
	if err != nil {
		log.Panicf("failed to initialize instrument: %v", err)
	}

	counter.Add(ctx, 1, attribute.KeyValue{
		Key:   semconv.RPCGRPCStatusCodeKey,
		Value: attribute.StringValue(status.Code(err).String()),
	})
}
