package contracts

import "time"

// GRPCConfig объединяет конфигурацию gRPC сервера и клиента.
type GRPCConfig interface {
	// Настройки сервера
	Address() string
	Timeout() time.Duration
	IdleTimeout() time.Duration
	ShutdownTimeout() time.Duration

	// Настройки клиента
	MaxRecvMsgSize() int
	MaxSendMsgSize() int
	ClientTimeout() time.Duration
}
