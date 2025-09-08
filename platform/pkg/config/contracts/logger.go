package contracts

// LoggerModule содержит конфигурацию логирования.
// Вынесен в отдельный интерфейс для лучшего разделения ответственности
// и возможности переиспользования в других модулях.
type LoggerModule interface {
	Level() string
	AsJSON() bool
	OTLP() OTLPModule
}

// OTLPModule содержит конфигурацию OpenTelemetry для интеграции с OTLP коллектором.
type OTLPModule interface {
	Enable() bool
	Endpoint() string
	ShutdownTimeout() int
}
