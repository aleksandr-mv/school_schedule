package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// TimeoutInterceptor устанавливает таймаут выполнения unary gRPC-обработчика
// посредством контекста. По истечении таймаута контекст будет отменён.
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return handler(ctx, req)
	}
}
