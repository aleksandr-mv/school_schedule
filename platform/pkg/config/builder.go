package config

import (
	"fmt"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/app"
	grpcmodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/grpc"
	kafkamodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/kafka"
	loggermodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/logger"
	metricmodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/metric"
	mongomodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/mongo"
	postgresmodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/postgres"
	redismodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/redis"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/services"
	sessionmodule "github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/session"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/tracing"
)

// =============================================================================
// МОДУЛЬНАЯ СБОРКА КОНФИГУРАЦИИ
// =============================================================================

// buildModularConfig создает конфигурацию через модульную архитектуру.
// Каждый модуль загружается независимо и возвращает интерфейс.
// Это обеспечивает лучшую изоляцию и тестируемость компонентов.
func buildModularConfig() (*config, error) {
	// gRPC модуль
	grpcCfg, err := grpcmodule.New()
	if err != nil {
		return nil, fmt.Errorf("initializing grpc module: %w", err)
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

	// PostgreSQL модуль
	postgresCfg, err := postgresmodule.New()
	if err != nil {
		return nil, fmt.Errorf("initializing postgres module: %w", err)
	}

	// MongoDB модуль
	mongoCfg, err := mongomodule.New()
	if err != nil {
		return nil, fmt.Errorf("initializing mongo module: %w", err)
	}

	// Services модуль (внешние сервисы)
	servicesCfg, err := services.New()
	if err != nil {
		return nil, fmt.Errorf("initializing services module: %w", err)
	}

	// Kafka модуль (Apache Kafka)
	kafkaCfg, err := kafkamodule.New()
	if err != nil {
		return nil, fmt.Errorf("initializing kafka module: %w", err)
	}

	// Redis модуль (Redis кэш и сессии)
	redisCfg, err := redismodule.New()
	if err != nil {
		return nil, fmt.Errorf("initializing redis module: %w", err)
	}

	// Session модуль (конфигурация сессий)
	sessionCfg, err := sessionmodule.New()
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
		grpcConfig:     grpcCfg,
		appConfig:      appCfg,
		loggerConfig:   loggerCfg,
		postgresConfig: postgresCfg,
		mongoConfig:    mongoCfg,
		servicesConfig: servicesCfg,
		kafkaConfig:    kafkaCfg,
		redisConfig:    redisCfg,
		sessionConfig:  sessionCfg,
		tracingConfig:  tracingCfg,
		metricConfig:   metricCfg,
	}, nil
}
