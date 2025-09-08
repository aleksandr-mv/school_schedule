package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// TimeoutUnaryClientInterceptor устанавливает таймаут на каждый исходящий unary RPC.
func TimeoutUnaryClientInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		cctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return invoker(cctx, method, req, reply, cc, opts...)
	}
}
