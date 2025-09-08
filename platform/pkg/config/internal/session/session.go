package session

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// rawSessionConfig соответствует секции session в YAML и env-переменным
type rawSessionConfig struct {
	TTL time.Duration `mapstructure:"ttl" yaml:"ttl" env:"SESSION_TTL"`
}

// sessionConfig хранит данные из секции session и реализует SessionConfig
type sessionConfig struct {
	Raw rawSessionConfig
}

var _ contracts.SessionConfig = (*sessionConfig)(nil)

// defaultSessionConfig возвращает конфигурацию сессий с дефолтными значениями
func defaultSessionConfig() *rawSessionConfig {
	return &rawSessionConfig{
		TTL: 24 * time.Hour,
	}
}

// DefaultSessionConfig читает конфигурацию сессий из ENV
func DefaultSessionConfig() (*sessionConfig, error) {
	raw := defaultSessionConfig()
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse session env: %w", err)
	}
	return &sessionConfig{Raw: *raw}, nil
}

// NewSessionConfig создает конфигурацию сессий. YAML -> ENV -> валидация
func NewSessionConfig() (*sessionConfig, error) {
	if section := helpers.GetSection("session"); section != nil {
		raw := defaultSessionConfig()
		if err := section.Unmarshal(raw); err == nil {
			return &sessionConfig{Raw: *raw}, nil
		}
	}
	return DefaultSessionConfig()
}

func (s *sessionConfig) TTL() time.Duration {
	return s.Raw.TTL
}
