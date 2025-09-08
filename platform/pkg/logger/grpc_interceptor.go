package logger

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor создает gRPC unary interceptor для логирования запросов.
// Интерсептор логирует начало/конец gRPC-запроса с минимально необходимыми полями.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		Info(ctx, " [gRPC] Запрос начат", zap.String("method", info.FullMethod))

		resp, err := handler(ctx, req)

		if err != nil {
			st, _ := status.FromError(err)
			Error(ctx, "❌ [gRPC] Запрос завершился ошибкой",
				zap.String("method", info.FullMethod),
				zap.String("grpc_code", st.Code().String()),
				zap.Error(err),
			)
		} else {
			Info(ctx, "✅ [gRPC] Запрос завершён", zap.String("method", info.FullMethod))
		}

		return resp, err
	}
}
