package logger

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// rawLoggerConfig соответствует полям секции logger в YAML и env-переменным.
type rawLoggerConfig struct {
	// Логирование
	Level  string `mapstructure:"level"   yaml:"level"   env:"LOG_LEVEL"`
	AsJSON bool   `mapstructure:"as_json" yaml:"as_json" env:"LOG_AS_JSON"`

	// OpenTelemetry интеграция
	OTLP rawOTLPConfig `mapstructure:"otlp" yaml:"otlp"`
}

// loggerConfig хранит данные секции logger и реализует LoggerConfig.
type loggerConfig struct {
	Raw rawLoggerConfig `yaml:"logger"`
}

// defaultLoggerConfig возвращает конфигурацию логгера с дефолтными значениями
func defaultLoggerConfig() *rawLoggerConfig {
	return &rawLoggerConfig{
		Level:  "info",
		AsJSON: true,
		OTLP:   *defaultOTLPConfig(),
	}
}

// DefaultLoggerConfig читает конфигурацию логгера из ENV.
func DefaultLoggerConfig() (*loggerConfig, error) {
	raw := defaultLoggerConfig()

	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse logger env: %w", err)
	}

	return &loggerConfig{Raw: *raw}, nil
}

// NewLoggerConfig создает конфигурацию логгера. YAML -> ENV -> валидация.
func NewLoggerConfig() (*loggerConfig, error) {
	if section := helpers.GetSection("logger"); section != nil {
		raw := defaultLoggerConfig()

		// Применяем YAML конфигурацию поверх дефолтов
		if err := section.Unmarshal(raw); err == nil {
			return &loggerConfig{Raw: *raw}, nil
		}
	}

	return DefaultLoggerConfig()
}

// Методы логирования
func (c *loggerConfig) Level() string { return c.Raw.Level }
func (c *loggerConfig) AsJSON() bool  { return c.Raw.AsJSON }

// OTLP возвращает модуль OTLP конфигурации
func (c *loggerConfig) OTLP() contracts.OTLPModule {
	return NewOTLPModule(&c.Raw.OTLP)
}
