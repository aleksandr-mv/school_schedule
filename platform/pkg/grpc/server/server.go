package grpcserver

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/grpc/interceptor"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
)

// BuildUnaryInterceptors строит базовую цепочку Unary-интерсепторов сервера gRPC
// на основе конфигурации платформы (логирование, recovery, валидация, таймаут).
func BuildUnaryInterceptors(timeout time.Duration) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		logger.UnaryServerInterceptor(),
		interceptor.RecoveryInterceptor(),
		interceptor.ValidationInterceptor(),
		interceptor.TimeoutInterceptor(timeout),
	}
}

// New создаёт gRPC-сервер с базовыми интерсепторами и опциональными доп. интерсепторами.
func New(_ context.Context,
	timeout time.Duration,
	maxRecvMsgSize, maxSendMsgSize int,
	extra ...grpc.UnaryServerInterceptor,
) *grpc.Server {
	chain := append(BuildUnaryInterceptors(timeout), extra...)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(chain...)),
		grpc.MaxRecvMsgSize(maxRecvMsgSize),
		grpc.MaxSendMsgSize(maxSendMsgSize),
	}

	return grpc.NewServer(opts...)
}
