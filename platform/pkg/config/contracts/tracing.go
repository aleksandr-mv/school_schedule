package contracts

import "time"

// TracingConfig содержит конфигурацию трейсинга OpenTelemetry.
type TracingConfig interface {
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
