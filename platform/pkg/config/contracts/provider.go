// Package contracts содержит все интерфейсы модульной конфигурации.
//
// Архитектура основана на принципах:
// - Модульность: каждый домен в отдельном интерфейсе
// - Инкапсуляция: внутренние типы скрыты от внешнего мира
// - Полиморфизм: интерфейсы как контракты поведения
// - Композиция: Provider объединяет модули
//
// Основные принципы:
//   - Интерфейсы описывают только чтение конфигурации (getters)
//   - Immutable доступ к данным (возвращаются копии для slices)
//   - Модульная организация по доменам
//   - Обратная совместимость через делегирование
package contracts

// ============================================================================
// ГЛАВНЫЙ ПРОВАЙДЕР КОНФИГУРАЦИИ
// ============================================================================

// Provider обеспечивает единый доступ к модульной конфигурации приложения.
//
// Чистая модульная архитектура основана на принципах:
// - Разделение ответственности: каждый модуль отвечает за свой домен
// - Инкапсуляция: внутренние детали скрыты за интерфейсами
// - Композиция: Provider объединяет независимые модули
// - Простота: минимальный и понятный API
//
// Доступные модули:
//   - Network(): сетевые настройки (CORS, TLS)
//   - Transport(): транспортные протоколы (HTTP, gRPC)
//   - App(): приложение и логирование
//   - Database(): агрегация баз данных (PostgreSQL, MongoDB)
//   - Redis(): конфигурация Redis кэша
//   - Services(): внешние сервисы
//   - Kafka(): конфигурация Apache Kafka
//
// Пример использования:
//
//	cfg := config.Load(ctx)
//	httpServer := cfg.Transport().HTTP()
//	corsConfig := cfg.Network().CORS()
//	loggerLevel := cfg.Logger().LoggerLevel()
//
//	// Kafka опциональна
//	if cfg.Kafka().IsEnabled() {
//		kafkaBrokers := cfg.Kafka().Brokers()
//	}
type Provider interface {
	// Network возвращает модуль сетевых настроек
	Network() NetworkConfig

	// Transport возвращает модуль транспортных протоколов
	Transport() TransportModule

	// App возвращает модуль приложения
	App() AppModule

	// Logger возвращает модуль логирования
	Logger() LoggerModule

	// Database возвращает агрегированную конфигурацию баз данных
	Database() DatabaseConfig

	// Services возвращает конфигурацию внешних сервисов
	Services() ServicesConfig

	// Redis возвращает конфигурацию Redis кэша
	Redis() RedisConfig

	// Kafka возвращает конфигурацию Apache Kafka
	Kafka() KafkaConfig

	// Telegram возвращает конфигурацию Telegram Bot
	Telegram() TelegramConfig

	// Session возвращает конфигурацию сессий
	Session() SessionConfig

	// Tracing возвращает конфигурацию трейсинга OpenTelemetry
	Tracing() TracingModule

	// Metric возвращает конфигурацию метрик OpenTelemetry
	Metric() MetricModule
}
