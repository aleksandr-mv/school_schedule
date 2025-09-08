package app

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// rawAppConfig соответствует полям секции app в YAML и env-переменным.
type rawAppConfig struct {
	ServiceName        string `mapstructure:"name" yaml:"name" env:"NAME"`
	ServiceEnvironment string `mapstructure:"environment" yaml:"environment" env:"ENVIRONMENT"`
	ServiceVersion     string `mapstructure:"version" yaml:"version" env:"VERSION"`
	MigrationsDir      string `mapstructure:"migrations_dir" yaml:"migrations_dir" env:"MIGRATIONS_DIR"`
	SwaggerPath        string `mapstructure:"swagger_path" yaml:"swagger_path" env:"SWAGGER_PATH"`
	SwaggerUIPath      string `mapstructure:"swagger_ui_path" yaml:"swagger_ui_path" env:"SWAGGER_UI_PATH"`
}

// appConfig хранит данные секции app и реализует AppConfig.
type appConfig struct {
	Raw rawAppConfig `yaml:"app"`
}

// defaultAppConfig возвращает конфигурацию приложения с дефолтными значениями
func defaultAppConfig() *rawAppConfig {
	return &rawAppConfig{
		ServiceName:        "unknown-service",
		ServiceEnvironment: "development",
		ServiceVersion:     "1.0.0",
		MigrationsDir:      "",
		SwaggerPath:        "",
		SwaggerUIPath:      "",
	}
}

// DefaultAppConfig читает конфигурацию приложения из ENV.
func DefaultAppConfig() (*appConfig, error) {
	raw := defaultAppConfig()
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse app env: %w", err)
	}
	return &appConfig{Raw: *raw}, nil
}

// NewAppConfig создает конфигурацию приложения. YAML -> ENV -> валидация.
func NewAppConfig() (*appConfig, error) {
	if section := helpers.GetSection("app"); section != nil {
		raw := defaultAppConfig()
		if err := section.Unmarshal(raw); err == nil {
			return &appConfig{Raw: *raw}, nil
		}
	}
	return DefaultAppConfig()
}

func (c *appConfig) Name() string          { return c.Raw.ServiceName }
func (c *appConfig) Environment() string   { return c.Raw.ServiceEnvironment }
func (c *appConfig) Version() string       { return c.Raw.ServiceVersion }
func (c *appConfig) MigrationsDir() string { return c.Raw.MigrationsDir }
func (c *appConfig) SwaggerPath() string   { return c.Raw.SwaggerPath }
func (c *appConfig) SwaggerUIPath() string { return c.Raw.SwaggerUIPath }
