package redis

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// rawRedisPoolConfig соответствует секции redis.pool в YAML и env-переменным.
type rawRedisPoolConfig struct {
	PoolSize      int           `mapstructure:"pool_size"       yaml:"pool_size"       env:"REDIS_POOL_SIZE"`
	MinIdle       int           `mapstructure:"min_idle"        yaml:"min_idle"        env:"REDIS_MIN_IDLE"`
	MaxRetries    int           `mapstructure:"max_retries"     yaml:"max_retries"     env:"REDIS_MAX_RETRIES"`
	DialTimeout   time.Duration `mapstructure:"dial_timeout"    yaml:"dial_timeout"    env:"REDIS_DIAL_TIMEOUT"`
	ReadTimeout   time.Duration `mapstructure:"read_timeout"    yaml:"read_timeout"    env:"REDIS_READ_TIMEOUT"`
	WriteTimeout  time.Duration `mapstructure:"write_timeout"   yaml:"write_timeout"   env:"REDIS_WRITE_TIMEOUT"`
	PoolTimeout   time.Duration `mapstructure:"pool_timeout"    yaml:"pool_timeout"    env:"REDIS_POOL_TIMEOUT"`
	IdleTimeout   time.Duration `mapstructure:"idle_timeout"    yaml:"idle_timeout"    env:"REDIS_IDLE_TIMEOUT"`
	IdleCheckFreq time.Duration `mapstructure:"idle_check_freq" yaml:"idle_check_freq" env:"REDIS_IDLE_CHECK_FREQ"`
}

// redisPoolConfig хранит данные из секции redis.pool и реализует RedisPoolConfig.
type redisPoolConfig struct {
	Raw rawRedisPoolConfig
}

var _ contracts.RedisPoolConfig = (*redisPoolConfig)(nil)

// defaultRedisPoolConfig возвращает конфигурацию Redis Pool с дефолтными значениями
func defaultRedisPoolConfig() *rawRedisPoolConfig {
	return &rawRedisPoolConfig{
		PoolSize:      10,
		MinIdle:       5,
		MaxRetries:    3,
		DialTimeout:   5 * time.Second,
		ReadTimeout:   3 * time.Second,
		WriteTimeout:  3 * time.Second,
		PoolTimeout:   4 * time.Second,
		IdleTimeout:   5 * time.Minute,
		IdleCheckFreq: 1 * time.Minute,
	}
}

// DefaultRedisPoolConfig читает Redis Pool из ENV.
func DefaultRedisPoolConfig() (*redisPoolConfig, error) {
	// Начинаем с дефолтной конфигурации
	raw := defaultRedisPoolConfig()

	// Применяем переменные окружения поверх дефолтов
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse redis pool env: %w", err)
	}

	return &redisPoolConfig{Raw: *raw}, nil
}

// NewRedisPoolConfig создает конфигурацию Redis Pool. YAML -> ENV -> валидация.
func NewRedisPoolConfig() (*redisPoolConfig, error) {
	if section := helpers.GetSection("redis.pool"); section != nil {
		// Начинаем с дефолтной конфигурации
		raw := defaultRedisPoolConfig()

		// Применяем YAML конфигурацию поверх дефолтов
		if err := section.Unmarshal(raw); err == nil {
			cfg := &redisPoolConfig{Raw: *raw}
			if err = cfg.validate(); err != nil {
				return nil, err
			}
			return cfg, nil
		}
	}

	return DefaultRedisPoolConfig()
}

func (r *redisPoolConfig) PoolSize() int {
	return r.Raw.PoolSize
}

func (r *redisPoolConfig) MinIdle() int {
	return r.Raw.MinIdle
}

func (r *redisPoolConfig) MaxRetries() int {
	return r.Raw.MaxRetries
}

func (r *redisPoolConfig) DialTimeout() time.Duration {
	return r.Raw.DialTimeout
}

func (r *redisPoolConfig) ReadTimeout() time.Duration {
	return r.Raw.ReadTimeout
}

func (r *redisPoolConfig) WriteTimeout() time.Duration {
	return r.Raw.WriteTimeout
}

func (r *redisPoolConfig) PoolTimeout() time.Duration {
	return r.Raw.PoolTimeout
}

func (r *redisPoolConfig) IdleTimeout() time.Duration {
	return r.Raw.IdleTimeout
}

func (r *redisPoolConfig) IdleCheckFreq() time.Duration {
	return r.Raw.IdleCheckFreq
}

func (r *redisPoolConfig) String() string {
	return fmt.Sprintf(
		"RedisPool{PoolSize:%d, MinIdle:%d, MaxRetries:%d, DialTimeout:%v, ReadTimeout:%v, WriteTimeout:%v, PoolTimeout:%v, IdleTimeout:%v, IdleCheckFreq:%v}",
		r.Raw.PoolSize, r.Raw.MinIdle, r.Raw.MaxRetries, r.Raw.DialTimeout, r.Raw.ReadTimeout, r.Raw.WriteTimeout, r.Raw.PoolTimeout, r.Raw.IdleTimeout, r.Raw.IdleCheckFreq,
	)
}

func (r *redisPoolConfig) validate() error {
	if r.Raw.PoolSize <= 0 {
		return fmt.Errorf("redis pool_size must be positive, got %d", r.Raw.PoolSize)
	}

	if r.Raw.MinIdle < 0 {
		return fmt.Errorf("redis min_idle cannot be negative, got %d", r.Raw.MinIdle)
	}

	if r.Raw.MinIdle > r.Raw.PoolSize {
		return fmt.Errorf("redis min_idle cannot be greater than pool_size, got %d > %d", r.Raw.MinIdle, r.Raw.PoolSize)
	}

	if r.Raw.MaxRetries < 0 {
		return fmt.Errorf("redis max_retries cannot be negative, got %d", r.Raw.MaxRetries)
	}

	if r.Raw.DialTimeout <= 0 {
		return fmt.Errorf("redis dial_timeout must be positive, got %v", r.Raw.DialTimeout)
	}

	if r.Raw.ReadTimeout <= 0 {
		return fmt.Errorf("redis read_timeout must be positive, got %v", r.Raw.ReadTimeout)
	}

	if r.Raw.WriteTimeout <= 0 {
		return fmt.Errorf("redis write_timeout must be positive, got %v", r.Raw.WriteTimeout)
	}

	if r.Raw.PoolTimeout <= 0 {
		return fmt.Errorf("redis pool_timeout must be positive, got %v", r.Raw.PoolTimeout)
	}

	if r.Raw.IdleTimeout <= 0 {
		return fmt.Errorf("redis idle_timeout must be positive, got %v", r.Raw.IdleTimeout)
	}

	if r.Raw.IdleCheckFreq <= 0 {
		return fmt.Errorf("redis idle_check_freq must be positive, got %v", r.Raw.IdleCheckFreq)
	}

	return nil
}
