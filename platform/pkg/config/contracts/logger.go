package contracts

// LoggerConfig содержит конфигурацию логирования.
// Вынесен в отдельный интерфейс для лучшего разделения ответственности
// и возможности переиспользования в других модулях.
type LoggerConfig interface {
	Level() string
	AsJSON() bool
	Enable() bool
	Endpoint() string
	ShutdownTimeout() int
}
