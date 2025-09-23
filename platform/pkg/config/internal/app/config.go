package app

import (
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
)

// rawConfig для загрузки данных из YAML/ENV
type rawConfig struct {
	ServiceName        string `mapstructure:"name" yaml:"name" env:"APP_NAME"`
	ServiceEnvironment string `mapstructure:"environment" yaml:"environment" env:"APP_ENVIRONMENT"`
	ServiceVersion     string `mapstructure:"version" yaml:"version" env:"APP_VERSION"`
	MigrationsDir      string `mapstructure:"migrations_dir" yaml:"migrations_dir" env:"MIGRATIONS_DIR"`
	SwaggerPath        string `mapstructure:"swagger_path" yaml:"swagger_path" env:"SWAGGER_PATH"`
	SwaggerUIPath      string `mapstructure:"swagger_ui_path" yaml:"swagger_ui_path" env:"SWAGGER_UI_PATH"`
}

// Config публичная структура для использования
type Config struct {
	raw rawConfig
}

// Компиляционная проверка
var _ contracts.AppConfig = (*Config)(nil)

// defaultConfig возвращает rawConfig с дефолтными значениями
func defaultConfig() rawConfig {
	return rawConfig{
		ServiceName:        "unknown-service",
		ServiceEnvironment: "development",
		ServiceVersion:     "1.0.0",
		MigrationsDir:      "",
		SwaggerPath:        "",
		SwaggerUIPath:      "",
	}
}

// Методы конфигурации
func (c *Config) Name() string          { return c.raw.ServiceName }
func (c *Config) Environment() string   { return c.raw.ServiceEnvironment }
func (c *Config) Version() string       { return c.raw.ServiceVersion }
func (c *Config) MigrationsDir() string { return c.raw.MigrationsDir }
func (c *Config) SwaggerPath() string   { return c.raw.SwaggerPath }
func (c *Config) SwaggerUIPath() string { return c.raw.SwaggerUIPath }
