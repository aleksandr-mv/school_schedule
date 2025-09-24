package grpcclient

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	grpcint "github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/grpc/interceptor"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
)

const (
	// TraceIDMetadataKey ключ для передачи trace ID в gRPC metadata
	TraceIDMetadataKey = "x-trace-id"
	// RequestIDMetadataKey ключ для передачи request ID в gRPC metadata
	RequestIDMetadataKey = "x-request-id"
)

// PropagateIDsUnaryClientInterceptor добавляет в исходящие gRPC metadata
// заголовки x-trace-id и x-request-id, извлечённые из контекста.
func PropagateIDsUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if traceID := logger.TraceIDFrom(ctx); traceID != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, TraceIDMetadataKey, traceID)
		}

		if requestID := logger.RequestIDFrom(ctx); requestID != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, RequestIDMetadataKey, requestID)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// PropagateSessionIDUnaryClientInterceptor добавляет session-id в исходящие gRPC metadata
func PropagateSessionIDUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// NewClient создаёт gRPC‑клиент (*grpc.ClientConn) с клиентскими интерсепторами по умолчанию,
// а также с дополнительными интерсепторами, переданными в extra.
func NewClient(address string, maxRecvMsgSize, maxSendMsgSize int, timeout time.Duration, extra ...grpc.UnaryClientInterceptor) (*grpc.ClientConn, error) {
	chain := append(
		[]grpc.UnaryClientInterceptor{
			PropagateIDsUnaryClientInterceptor(),
			PropagateSessionIDUnaryClientInterceptor(),
			grpcint.TimeoutUnaryClientInterceptor(timeout),
		},
		extra...,
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(chain...),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxRecvMsgSize),
			grpc.MaxCallSendMsgSize(maxSendMsgSize),
		),
	}

	return grpc.NewClient(address, opts...)
}
