package contracts

import "time"

// MetricModule содержит конфигурацию метрик OpenTelemetry.
type MetricConfig interface {
	// Основные настройки
	Enable() bool
	Endpoint() string
	Timeout() time.Duration

	// Настройки метрик
	Namespace() string
	AppName() string

	// Настройки экспорта
	ExportInterval() time.Duration
	ShutdownTimeout() time.Duration

	// Настройки bucket'ов для гистограмм
	BucketBoundaries() []float64
}
