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
// - Network:   CORS + TLS конфигурации
// - Transport: HTTP + gRPC серверы и клиенты
// - App:       Логирование + приложение
// - Database:  PostgreSQL + MongoDB
// - Services:  Внешние сервисы
// - Kafka:     Apache Kafka брокеры и топики
//
// Использование (API остался без изменений):
//
//	cfg, err := config.Load(ctx)
//	if err != nil { panic(err) }
//	httpServer := cfg.HTTP()     // Работает как раньше
//	dbConfig := cfg.DatabaseConfig()
package config

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/app"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/database"
	kafkamodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/kafka"
	loggermodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/logger"
	metricmodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/metric"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/network"
	redismodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/redis"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/services"
	sessionmodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/session"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/telegram"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/tracing"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/transport"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

// =============================================================================
// ОСНОВНАЯ СТРУКТУРА КОНФИГУРАЦИИ
// =============================================================================

// config реализует интерфейс contracts.Provider через модульную архитектуру.
// Каждое поле представляет отдельный модуль с собственной областью ответственности.
// Использует интерфейсы для обеспечения инкапсуляции и возможности замены реализаций.
type config struct {
	networkConfig   contracts.NetworkConfig   // Сетевые настройки (CORS, TLS)
	transportConfig contracts.TransportModule // Транспортные протоколы (HTTP, gRPC)
	appModule       contracts.AppModule       // Приложение
	loggerModule    contracts.LoggerModule    // Логирование
	dbConfig        contracts.DatabaseConfig  // Базы данных (PostgreSQL, MongoDB)
	servicesConfig  contracts.ServicesConfig  // Внешние сервисы
	kafkaConfig     contracts.KafkaConfig     // Apache Kafka
	telegramConfig  contracts.TelegramConfig  // Telegram Bot
	redisConfig     contracts.RedisConfig     // Redis кэш и сессии
	sessionConfig   contracts.SessionConfig   // Конфигурация сессий
	tracingModule   contracts.TracingModule   // Трейсинг OpenTelemetry
	metricModule    contracts.MetricModule    // Метрики OpenTelemetry
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
// ПРИВАТНЫЕ ФУНКЦИИ ЗАГРУЗКИ МОДУЛЬНОЙ КОНФИГУРАЦИИ
// =============================================================================

// buildModularConfig создает конфигурацию через модульную архитектуру.
// Каждый модуль загружается независимо и возвращает интерфейс.
// Это обеспечивает лучшую изоляцию и тестируемость компонентов.
func buildModularConfig() (*config, error) {
	// Network модуль (CORS + будущий TLS)
	networkCfg, err := network.New()
	if err != nil {
		return nil, fmt.Errorf("initializing network module: %w", err)
	}

	// Transport модуль (HTTP + gRPC серверы и клиенты)
	transportCfg, err := transport.New()
	if err != nil {
		return nil, fmt.Errorf("initializing transport module: %w", err)
	}

	// App модуль (приложение)
	appCfg, err := app.New()
	if err != nil {
		return nil, fmt.Errorf("initializing app module: %w", err)
	}

	// Logger модуль (логирование)
	loggerCfg, err := loggermodule.New()
	if err != nil {
		return nil, fmt.Errorf("initializing logger module: %w", err)
	}

	// Database модуль (PostgreSQL + MongoDB агрегация)
	dbCfg, err := database.New()
	if err != nil {
		return nil, fmt.Errorf("initializing database module: %w", err)
	}

	// Services модуль (внешние сервисы)
	servicesCfg, err := services.New()
	if err != nil {
		return nil, fmt.Errorf("initializing services module: %w", err)
	}

	// Kafka модуль (Apache Kafka)
	kafkaCfg, err := kafkamodule.NewKafkaConfig()
	if err != nil {
		return nil, fmt.Errorf("initializing kafka module: %w", err)
	}

	// Telegram модуль (Telegram Bot)
	telegramCfg, err := telegram.New()
	if err != nil {
		return nil, fmt.Errorf("initializing telegram module: %w", err)
	}

	// Redis модуль (Redis кэш и сессии)
	redisCfg, err := redismodule.New()
	if err != nil {
		return nil, fmt.Errorf("initializing redis module: %w", err)
	}

	// Session модуль (конфигурация сессий)
	sessionCfg, err := sessionmodule.NewSessionConfig()
	if err != nil {
		return nil, fmt.Errorf("initializing session module: %w", err)
	}

	// Tracing модуль (OpenTelemetry трейсинг)
	tracingCfg, err := tracing.New()
	if err != nil {
		return nil, fmt.Errorf("initializing tracing module: %w", err)
	}

	// Metric модуль (OpenTelemetry метрики)
	metricCfg, err := metricmodule.New()
	if err != nil {
		return nil, fmt.Errorf("initializing metric module: %w", err)
	}

	return &config{
		networkConfig:   networkCfg,
		transportConfig: transportCfg,
		appModule:       appCfg,
		loggerModule:    loggerCfg,
		dbConfig:        dbCfg,
		servicesConfig:  servicesCfg,
		kafkaConfig:     kafkaCfg,
		telegramConfig:  telegramCfg,
		redisConfig:     redisCfg,
		sessionConfig:   sessionCfg,
		tracingModule:   tracingCfg,
		metricModule:    metricCfg,
	}, nil
}

// =============================================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА contracts.Provider
// =============================================================================

// Реализация чистого модульного Provider интерфейса.
// Каждый метод возвращает соответствующий модуль без дополнительных оберток.
//
// Использование:
//   cfg := config.Load(ctx)
//   httpConfig := cfg.Transport().HTTP()
//   corsConfig := cfg.Network().CORS()
//   loggerLevel := cfg.Logger().LoggerLevel()
//   dbConfig := cfg.Database()
//
//   // Kafka опциональна - проверяем перед использованием
//   if cfg.Kafka().IsEnabled() {
//       kafkaBrokers := cfg.Kafka().Brokers()
//   }

// =============================================================================
// МОДУЛЬНАЯ РЕАЛИЗАЦИЯ PROVIDER ИНТЕРФЕЙСА
// =============================================================================

// Network возвращает модуль сетевых настроек (CORS, TLS)
func (c *config) Network() contracts.NetworkConfig {
	return c.networkConfig
}

// Transport возвращает модуль транспортных протоколов (HTTP, gRPC)
func (c *config) Transport() contracts.TransportModule {
	return c.transportConfig
}

// App возвращает модуль приложения
func (c *config) App() contracts.AppModule {
	return c.appModule
}

// Logger возвращает модуль логирования
func (c *config) Logger() contracts.LoggerModule {
	return c.loggerModule
}

// Database возвращает агрегированную конфигурацию баз данных
func (c *config) Database() contracts.DatabaseConfig {
	return c.dbConfig
}

// Services возвращает конфигурацию внешних сервисов
func (c *config) Services() contracts.ServicesConfig {
	return c.servicesConfig
}

// Kafka возвращает конфигурацию Apache Kafka
func (c *config) Kafka() contracts.KafkaConfig {
	return c.kafkaConfig
}

// Telegram возвращает конфигурацию Telegram Bot
func (c *config) Telegram() contracts.TelegramConfig {
	return c.telegramConfig
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
func (c *config) Tracing() contracts.TracingModule {
	return c.tracingModule
}

// Metric возвращает конфигурацию метрик OpenTelemetry
func (c *config) Metric() contracts.MetricModule {
	return c.metricModule
}
