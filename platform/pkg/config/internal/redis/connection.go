package redis

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// rawRedisConnectionConfig соответствует секции redis.connection в YAML и env-переменным.
type rawRedisConnectionConfig struct {
	Host     string `mapstructure:"host"     yaml:"host"     env:"REDIS_HOST"`
	Port     int    `mapstructure:"port"     yaml:"port"     env:"REDIS_PORT"`
	Password string `mapstructure:"password" yaml:"password" env:"REDIS_PASSWORD"`
	Database int    `mapstructure:"database" yaml:"database" env:"REDIS_DATABASE"`
}

// redisConnectionConfig хранит данные из секции redis.connection и реализует RedisConnection.
type redisConnectionConfig struct {
	Raw rawRedisConnectionConfig
}

var _ contracts.RedisConnection = (*redisConnectionConfig)(nil)

// defaultRedisConnectionConfig возвращает конфигурацию Redis Connection с дефолтными значениями
func defaultRedisConnectionConfig() *rawRedisConnectionConfig {
	return &rawRedisConnectionConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Database: 0,
	}
}

// DefaultRedisConnectionConfig читает Redis Connection из ENV.
func DefaultRedisConnectionConfig() (*redisConnectionConfig, error) {
	// Начинаем с дефолтной конфигурации
	raw := defaultRedisConnectionConfig()

	// Применяем переменные окружения поверх дефолтов
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse redis connection env: %w", err)
	}

	return &redisConnectionConfig{Raw: *raw}, nil
}

// NewRedisConnectionConfig создает конфигурацию Redis Connection. YAML -> ENV -> валидация.
func NewRedisConnectionConfig() (*redisConnectionConfig, error) {
	if section := helpers.GetSection("redis.connection"); section != nil {
		// Начинаем с дефолтной конфигурации
		raw := defaultRedisConnectionConfig()

		// Применяем YAML конфигурацию поверх дефолтов
		if err := section.Unmarshal(raw); err == nil {
			return &redisConnectionConfig{Raw: *raw}, nil
		}
	}

	return DefaultRedisConnectionConfig()
}

func (r *redisConnectionConfig) Host() string {
	return r.Raw.Host
}

func (r *redisConnectionConfig) Port() int {
	return r.Raw.Port
}

func (r *redisConnectionConfig) Password() string {
	return r.Raw.Password
}

func (r *redisConnectionConfig) Database() int {
	return r.Raw.Database
}

func (r *redisConnectionConfig) Address() string {
	return fmt.Sprintf("%s:%d", r.Raw.Host, r.Raw.Port)
}

func (r *redisConnectionConfig) URI() string {
	if r.Raw.Password != "" {
		return fmt.Sprintf("redis://:%s@%s:%d/%d", r.Raw.Password, r.Raw.Host, r.Raw.Port, r.Raw.Database)
	}
	return fmt.Sprintf("redis://%s:%d/%d", r.Raw.Host, r.Raw.Port, r.Raw.Database)
}

func (r *redisConnectionConfig) String() string {
	return fmt.Sprintf(
		"RedisConnection{Host:%s, Port:%d, Database:%d}",
		r.Raw.Host, r.Raw.Port, r.Raw.Database,
	)
}
