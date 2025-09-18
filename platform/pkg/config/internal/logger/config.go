package logger

import (
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.LoggerConfig = (*Config)(nil)

// rawConfig для загрузки данных из YAML/ENV
type rawConfig struct {
	// Логирование
	Level  string `mapstructure:"level"   yaml:"level"   env:"LOG_LEVEL"`
	AsJSON bool   `mapstructure:"as_json" yaml:"as_json" env:"LOG_AS_JSON"`

	// OpenTelemetry интеграция
	Enable          bool   `mapstructure:"enable" yaml:"enable" env:"OTLP_ENABLE"`
	Endpoint        string `mapstructure:"endpoint" yaml:"endpoint" env:"OTLP_ENDPOINT"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout" env:"OTLP_SHUTDOWN_TIMEOUT"`
}

// Config публичная структура Logger конфигурации
type Config struct {
	raw rawConfig
}

// defaultConfig возвращает rawConfig с дефолтными значениями
func defaultConfig() rawConfig {
	return rawConfig{
		Level:           "info",
		AsJSON:          true,
		Enable:          false,
		Endpoint:        "localhost:4317",
		ShutdownTimeout: 2,
	}
}

// Методы для LoggerConfig интерфейса
func (c *Config) Level() string        { return c.raw.Level }
func (c *Config) AsJSON() bool         { return c.raw.AsJSON }
func (c *Config) Enable() bool         { return c.raw.Enable }
func (c *Config) Endpoint() string     { return c.raw.Endpoint }
func (c *Config) ShutdownTimeout() int { return c.raw.ShutdownTimeout }
