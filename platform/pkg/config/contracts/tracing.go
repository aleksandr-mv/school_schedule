package contracts

import "time"

// TracingModule содержит конфигурацию трейсинга OpenTelemetry.
type TracingModule interface {
	// Основные настройки
	Enable() bool
	Endpoint() string
	Timeout() time.Duration

	// Настройки семплирования
	SampleRatio() int

	// Настройки повторных попыток
	RetryEnabled() bool
	RetryInitialInterval() time.Duration
	RetryMaxInterval() time.Duration
	RetryMaxElapsedTime() time.Duration

	// Настройки пропагации
	EnableTraceContext() bool
	EnableBaggage() bool

	// Настройки shutdown
	ShutdownTimeout() time.Duration
}
