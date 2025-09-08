package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidationInterceptor выполняет валидацию входящего запроса, если он
// реализует интерфейс Validate() error. При ошибке возвращает codes.InvalidArgument.
func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if v, ok := req.(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
			}
		}
		return handler(ctx, req)
	}
}
