// Package config предоставляет модульную систему управления конфигурацией для микросервисов.
//
// Новая архитектура основана на принципах:
// - Модульность: каждый домен в отдельном пакете
// - Интерфейсы как контракты: полиморфизм и инкапсуляция
// - Типобезопасность: строгая типизация через contracts
// - Валидация: проверка всех параметров конфигурации
// - Переиспользование: существующий код loader'ов сохранен
//
// Архитектура модулей:
// - App:       Базовые настройки сервиса (имя, версия, окружение)
// - Logger:    Логирование с OpenTelemetry интеграцией
// - GRPC:      gRPC сервер и клиентские настройки
// - Postgres:  PostgreSQL Primary-Replica архитектура
// - Mongo:     MongoDB Primary-Replica архитектура
// - Redis:     Redis Cluster (3 шарда + 3 реплики)
// - Services:  Внешние микросервисы
// - Kafka:     Apache Kafka брокеры, consumers, producers
// - Telegram:  Telegram Bot интеграция
// - Session:   Конфигурация сессий
// - Tracing:   OpenTelemetry трейсинг
// - Metric:    OpenTelemetry метрики
//
// Использование:
//
//	cfg, err := config.Load(ctx)
//	if err != nil { panic(err) }
//
//	// Базовые модули
//	appName := cfg.App().Name()
//	logLevel := cfg.Logger().Level()
//	grpcAddr := cfg.GRPC().Address()
//
//	// Базы данных (Primary-Replica)
//	pgWrite := cfg.Postgres().PrimaryURI()
//	pgRead := cfg.Postgres().ReplicaURI()
//	mongoWrite := cfg.Mongo().PrimaryURI()
//
//	// Опциональные модули
//	if cfg.Redis().Cluster().IsEnabled() {
//		redisNodes := cfg.Redis().Cluster().Nodes()
//	}
//	if cfg.Kafka().IsEnabled() {
//		kafkaBrokers := cfg.Kafka().Brokers()
//	}
package config

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/helpers"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
)

// =============================================================================
// ОСНОВНАЯ СТРУКТУРА КОНФИГУРАЦИИ
// =============================================================================

// config реализует интерфейс contracts.Provider через модульную архитектуру.
// Каждое поле представляет отдельный модуль с собственной областью ответственности.
// Использует интерфейсы для обеспечения инкапсуляции и возможности замены реализаций.
type config struct {
	grpcConfig     contracts.GRPCConfig     // gRPC конфигурация
	appConfig      contracts.AppConfig      // Приложение
	loggerConfig   contracts.LoggerConfig   // Логирование
	postgresConfig contracts.PostgresConfig // PostgreSQL база данных
	mongoConfig    contracts.MongoConfig    // MongoDB база данных
	servicesConfig contracts.ServicesConfig // Внешние сервисы
	kafkaConfig    contracts.KafkaConfig    // Apache Kafka
	redisConfig    contracts.RedisConfig    // Redis кэш и сессии
	sessionConfig  contracts.SessionConfig  // Конфигурация сессий
	tracingConfig  contracts.TracingConfig  // Трейсинг OpenTelemetry
	metricConfig   contracts.MetricConfig   // Метрики OpenTelemetry
}

// =============================================================================
// ПУБЛИЧНЫЕ ФУНКЦИИ ИНИЦИАЛИЗАЦИИ
// =============================================================================

// Стратегия загрузки модульной конфигурации:
// 1. Инициализация Viper для чтения YAML и ENV переменных
// 2. Создание модулей через их конструкторы New()
// 3. Валидация каждого модуля независимо
// 4. Сборка единого Provider интерфейса
// 5. Обеспечение обратной совместимости с существующим API

// Load инициализирует модульную конфигурацию с сохранением обратной совместимости.
// Внешний API остается неизменным, но внутри используется новая модульная архитектура.
func Load(ctx context.Context) (contracts.Provider, error) {
	path := helpers.ResolveConfigPath("config/development.yaml")
	logger.Info(ctx, "🔍 [Config] Загружаем модульную конфигурацию", zap.String("path", path))

	if err := helpers.InitViper(path); err != nil {
		logger.Warn(ctx, "⚠️ [Config] Не удалось загрузить YAML, используем только ENV",
			zap.String("path", path), zap.Error(err))
	}

	cfg, err := buildModularConfig()
	if err != nil {
		return nil, err
	}

	logger.Info(ctx, "✅ [Config] Модульная конфигурация успешно загружена")
	return cfg, nil
}

// =============================================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА contracts.Provider
// =============================================================================

// Реализация чистого модульного Provider интерфейса.
// Каждый метод возвращает соответствующий модуль без дополнительных оберток.
//
// Использование:
//   cfg := config.Load(ctx)
//
//   // Используем модули
//   appName := cfg.App().Name()
//   grpcAddr := cfg.GRPC().Address()
//   pgURI := cfg.Postgres().PrimaryURI()
//
//   // Опциональные модули
//   if cfg.Redis().Cluster().IsEnabled() {
//       redisNodes := cfg.Redis().Cluster().Nodes()
//   }

// =============================================================================
// МОДУЛЬНАЯ РЕАЛИЗАЦИЯ PROVIDER ИНТЕРФЕЙСА
// =============================================================================

// App возвращает модуль приложения
func (c *config) App() contracts.AppConfig {
	return c.appConfig
}

// GRPC возвращает gRPC конфигурацию
func (c *config) GRPC() contracts.GRPCConfig {
	return c.grpcConfig
}

// Logger возвращает модуль логирования
func (c *config) Logger() contracts.LoggerConfig {
	return c.loggerConfig
}

// Database возвращает агрегированную конфигурацию баз данных
func (c *config) Postgres() contracts.PostgresConfig {
	return c.postgresConfig
}

func (c *config) Mongo() contracts.MongoConfig {
	return c.mongoConfig
}

// Services возвращает конфигурацию внешних сервисов
func (c *config) Services() contracts.ServicesConfig {
	return c.servicesConfig
}

// Kafka возвращает конфигурацию Apache Kafka
func (c *config) Kafka() contracts.KafkaConfig {
	return c.kafkaConfig
}

// Redis возвращает конфигурацию Redis
func (c *config) Redis() contracts.RedisConfig {
	return c.redisConfig
}

// Session возвращает конфигурацию сессий
func (c *config) Session() contracts.SessionConfig {
	return c.sessionConfig
}

// Tracing возвращает конфигурацию трейсинга OpenTelemetry
func (c *config) Tracing() contracts.TracingConfig {
	return c.tracingConfig
}

// Metric возвращает конфигурацию метрик OpenTelemetry
func (c *config) Metric() contracts.MetricConfig {
	return c.metricConfig
}
