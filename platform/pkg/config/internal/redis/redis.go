package redis

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// rawRedisConfig соответствует секции redis в YAML и env-переменным.
type rawRedisConfig struct {
	Enabled bool `mapstructure:"enabled" yaml:"enabled" env:"REDIS_ENABLED"`
}

// redisConfig хранит данные из секции redis и реализует RedisConfig.
type redisConfig struct {
	Raw  rawRedisConfig
	conn contracts.RedisConnection
	pool contracts.RedisPoolConfig
}

var _ contracts.RedisConfig = (*redisConfig)(nil)

// defaultRedisConfig возвращает конфигурацию Redis с дефолтными значениями
func defaultRedisConfig() *rawRedisConfig {
	return &rawRedisConfig{
		Enabled: false,
	}
}

// DefaultRedisConfig читает Redis из ENV.
func DefaultRedisConfig() (*redisConfig, error) {
	// Начинаем с дефолтной конфигурации
	raw := defaultRedisConfig()

	if err := env.Parse(raw); err != nil {
		return &redisConfig{
			Raw: *raw,
		}, nil
	}

	if !raw.Enabled {
		return &redisConfig{
			Raw: *raw,
		}, nil
	}

	conn, err := DefaultRedisConnectionConfig()
	if err != nil {
		return &redisConfig{
			Raw: *defaultRedisConfig(),
		}, nil
	}

	pool, err := DefaultRedisPoolConfig()
	if err != nil {
		return &redisConfig{
			Raw: *defaultRedisConfig(),
		}, nil
	}

	return &redisConfig{
		Raw:  *raw,
		conn: conn,
		pool: pool,
	}, nil
}

// NewRedisConfig создает конфигурацию Redis. YAML -> ENV -> валидация.
func NewRedisConfig() (*redisConfig, error) {
	section := helpers.GetSection("redis")
	if section == nil {
		return DefaultRedisConfig()
	}

	raw := defaultRedisConfig()

	if err := section.Unmarshal(raw); err != nil {
		return DefaultRedisConfig()
	}

	if !raw.Enabled {
		return &redisConfig{
			Raw: *raw,
		}, nil
	}

	conn, err := NewRedisConnectionConfig()
	if err != nil {
		return &redisConfig{
			Raw: *defaultRedisConfig(),
		}, nil
	}

	pool, err := NewRedisPoolConfig()
	if err != nil {
		return &redisConfig{
			Raw: *defaultRedisConfig(),
		}, nil
	}

	return &redisConfig{
		Raw:  *raw,
		conn: conn,
		pool: pool,
	}, nil
}

func (r *redisConfig) IsEnabled() bool {
	return r.Raw.Enabled
}

func (r *redisConfig) Connection() contracts.RedisConnection {
	if !r.Raw.Enabled {
		return nil
	}
	return r.conn
}

func (r *redisConfig) Pool() contracts.RedisPoolConfig {
	if !r.Raw.Enabled {
		return nil
	}
	return r.pool
}

func (r *redisConfig) String() string {
	return fmt.Sprintf(
		"Redis{Enabled:%v, Connection:%v, Pool:%v}",
		r.Raw.Enabled, r.conn, r.pool,
	)
}
