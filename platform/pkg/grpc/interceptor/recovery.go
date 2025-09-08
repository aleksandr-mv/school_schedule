package interceptor

import (
	"context"
	"runtime/debug"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

// RecoveryInterceptor перехватывает паники в unary gRPC-обработчиках и возвращает внутр. ошибку.
func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(ctx, "💥 [gRPC] Panic в обработчике",
					zap.Any("panic", r),
					zap.ByteString("stacktrace", debug.Stack()),
				)
				err = status.Errorf(13, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
