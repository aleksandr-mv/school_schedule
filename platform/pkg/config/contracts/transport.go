package contracts

import "time"

type TransportModule interface {
	HTTP() HTTPServer
	GRPC() GRPCServer
	GRPCClientLimits() GRPCClientLimits
	TLS() TLSConfig
}

// HTTPConfig интерфейс для конфигурации HTTP сервера.
type HTTPServer interface {
	Address() string
	ReadHeaderTimeout() time.Duration
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	IdleTimeout() time.Duration
	MaxHeaderBytes() int
	ShutdownTimeout() time.Duration
	HandlerTimeout() time.Duration
}

// GRPCConfig интерфейс для конфигурации gRPC-сервера.
type GRPCServer interface {
	Address() string
	Timeout() time.Duration
	IdleTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

type GRPCClientLimits interface {
	MaxRecvMsgSize() int
	MaxSendMsgSize() int
	Timeout() time.Duration
}
