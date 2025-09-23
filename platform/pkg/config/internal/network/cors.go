package network

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/helpers"
)

// rawCorsConfig соответствует секции cors в YAML и env-переменным.
type rawCorsConfig struct {
	Origins     []string `mapstructure:"allowed_origins"   yaml:"allowed_origins"   env:"CORS_ALLOWED_ORIGINS"`
	Methods     []string `mapstructure:"allowed_methods"   yaml:"allowed_methods"   env:"CORS_ALLOWED_METHODS"`
	Headers     []string `mapstructure:"allowed_headers"   yaml:"allowed_headers"   env:"CORS_ALLOWED_HEADERS"`
	Exposed     []string `mapstructure:"exposed_headers"   yaml:"exposed_headers"   env:"CORS_EXPOSED_HEADERS"`
	Credentials bool     `mapstructure:"allow_credentials" yaml:"allow_credentials" env:"CORS_ALLOW_CREDENTIALS"`
	Age         int      `mapstructure:"max_age"           yaml:"max_age"           env:"CORS_MAX_AGE"`
}

// corsConfig хранит данные из секции cors и реализует CORSConfig.
// Тег yaml:"cors" применяется сразу к полю raw.
type corsConfig struct {
	Raw rawCorsConfig `yaml:"cors"`
}

var _ contracts.CORSConfig = (*corsConfig)(nil)

// defaultCORSConfig возвращает конфигурацию CORS с дефолтными значениями
func defaultCORSConfig() *rawCorsConfig {
	return &rawCorsConfig{
		Origins:     []string{"*"},
		Methods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		Headers:     []string{"*"},
		Exposed:     []string{},
		Credentials: false,
		Age:         86400, // 24 часа
	}
}

// DefaultCORSConfig читает CORS из ENV.
func DefaultCORSConfig() (*corsConfig, error) {
	raw := defaultCORSConfig()
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse cors env: %w", err)
	}
	return &corsConfig{Raw: *raw}, nil
}

// NewCORSConfig создает конфигурацию CORS. YAML -> ENV -> валидация.
func NewCORSConfig() (*corsConfig, error) {
	if section := helpers.GetSection("cors"); section != nil {
		raw := defaultCORSConfig()
		if err := section.Unmarshal(raw); err == nil {
			return &corsConfig{Raw: *raw}, nil
		}
	}
	return DefaultCORSConfig()
}

func (c *corsConfig) AllowedOrigins() []string { return append([]string(nil), c.Raw.Origins...) }
func (c *corsConfig) AllowedMethods() []string { return append([]string(nil), c.Raw.Methods...) }
func (c *corsConfig) AllowedHeaders() []string { return append([]string(nil), c.Raw.Headers...) }
func (c *corsConfig) ExposedHeaders() []string { return append([]string(nil), c.Raw.Exposed...) }
func (c *corsConfig) AllowCredentials() bool   { return c.Raw.Credentials }
func (c *corsConfig) MaxAge() int              { return c.Raw.Age }
