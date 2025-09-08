package postgres

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/viper"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

type rawPoolConfig struct {
	MaxCons           int           `mapstructure:"max_cons"             yaml:"max_cons"             env:"DB_MAX_CONS"`
	MinCons           int           `mapstructure:"min_cons"             yaml:"min_cons"             env:"DB_MIN_CONS"`
	MaxConnLifetime   time.Duration `mapstructure:"max_conn_lifetime"    yaml:"max_conn_lifetime"    env:"DB_MAX_CONN_LIFETIME"`
	MaxConnIdleTime   time.Duration `mapstructure:"max_conn_idle_time"   yaml:"max_conn_idle_time"   env:"DB_MAX_CONN_IDLE_TIME"`
	HealthCheckPeriod time.Duration `mapstructure:"health_check_period"   yaml:"health_check_period"   env:"DB_HEALTH_CHECK_PERIOD"`
	ConnectTimeout    time.Duration `mapstructure:"connect_timeout"      yaml:"connect_timeout"      env:"DB_CONNECT_TIMEOUT"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"     yaml:"shutdown_timeout"     env:"DB_SHUTDOWN_TIMEOUT"`
}

// PoolConfig хранит настройки пула соединений к PostgreSQL.
type poolConfig struct {
	Raw rawPoolConfig `yaml:"pool"`
}

// Реализуем интерфейсы
var _ contracts.PostgresPoolConfig = (*poolConfig)(nil)

func defaultPoolConfig() *rawPoolConfig {
	return &rawPoolConfig{
		MaxCons:           10,
		MinCons:           2,
		MaxConnLifetime:   30 * time.Minute,
		MaxConnIdleTime:   5 * time.Minute,
		HealthCheckPeriod: 1 * time.Minute,
		ConnectTimeout:    10 * time.Second,
		ShutdownTimeout:   30 * time.Second,
	}
}

// DefaultPoolConfig читает конфиг пула из ENV
func DefaultPoolConfig() (*poolConfig, error) {
	raw := defaultPoolConfig()
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}
	return &poolConfig{Raw: *raw}, nil
}

// NewPoolConfig читает pool-конфиг из поддерева Viper (ожидается database.postgres.pool).
// Если sub == nil или секция отсутствует — fallback на переменные окружения через DefaultPoolConfig().
func NewPoolConfig(sub *viper.Viper) (*poolConfig, error) {
	if sub == nil {
		return DefaultPoolConfig()
	}
	pool := sub.Sub("pool")
	if pool == nil {
		return DefaultPoolConfig()
	}

	raw := defaultPoolConfig()
	if err := pool.Unmarshal(raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pool config: %w", err)
	}

	return &poolConfig{Raw: *raw}, nil
}

// Реализация интерфейса contracts.PoolConfig
func (p *poolConfig) MaxSize() int                   { return p.Raw.MaxCons }
func (p *poolConfig) MinSize() int                   { return p.Raw.MinCons }
func (p *poolConfig) MaxIdleTime() time.Duration     { return p.Raw.MaxConnIdleTime }
func (p *poolConfig) ConnectTimeout() time.Duration  { return p.Raw.ConnectTimeout }
func (p *poolConfig) ShutdownTimeout() time.Duration { return p.Raw.ShutdownTimeout }

// Реализация интерфейса contracts.PostgresPoolConfig
func (p *poolConfig) MaxConnLifetime() time.Duration   { return p.Raw.MaxConnLifetime }
func (p *poolConfig) HealthCheckPeriod() time.Duration { return p.Raw.HealthCheckPeriod }
