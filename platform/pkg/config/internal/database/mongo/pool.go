package mongo

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/viper"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

type rawPoolConfig struct {
	MaxPoolSize     int           `mapstructure:"max_pool_size"    yaml:"max_pool_size"       env:"MONGO_MAX_POOL_SIZE"`
	MinPoolSize     int           `mapstructure:"min_pool_size"    yaml:"min_pool_size"       env:"MONGO_MIN_POOL_SIZE"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time" yaml:"max_conn_idle_time"  env:"MONGO_MAX_CONN_IDLE_TIME"`
	MaxConnecting   int           `mapstructure:"max_connecting"    yaml:"max_connecting"      env:"MONGO_MAX_CONNECTING"`
	ConnectTimeout  time.Duration `mapstructure:"connect_timeout"   yaml:"connect_timeout"     env:"MONGO_CONNECT_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"  yaml:"shutdown_timeout"    env:"MONGO_SHUTDOWN_TIMEOUT"`
}

// PoolConfig хранит настройки пула соединений.
type poolConfig struct {
	Raw rawPoolConfig
}

// Реализуем интерфейсы
var _ contracts.MongoPoolConfig = (*poolConfig)(nil)

// defaultPoolConfig возвращает конфигурацию пула соединений MongoDB с дефолтными значениями
func defaultPoolConfig() *rawPoolConfig {
	return &rawPoolConfig{
		MaxPoolSize:     100,
		MinPoolSize:     5,
		MaxConnIdleTime: 30 * time.Minute,
		MaxConnecting:   2,
		ConnectTimeout:  10 * time.Second,
		ShutdownTimeout: 30 * time.Second,
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

// NewPoolConfig читает pool-конфиг Mongo из поддерева Viper (ожидается database.mongo.pool).
// При отсутствии секции — fallback на ENV.
func NewPoolConfig(sub *viper.Viper) (*poolConfig, error) {
	if sub == nil {
		return DefaultPoolConfig()
	}

	poolV := sub.Sub("pool")
	if poolV == nil {
		return DefaultPoolConfig()
	}

	raw := defaultPoolConfig()
	if err := poolV.Unmarshal(raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Mongo pool config: %w", err)
	}

	return &poolConfig{Raw: *raw}, nil
}

func (p *poolConfig) MaxSize() int                   { return p.Raw.MaxPoolSize }
func (p *poolConfig) MinSize() int                   { return p.Raw.MinPoolSize }
func (p *poolConfig) MaxIdleTime() time.Duration     { return p.Raw.MaxConnIdleTime }
func (p *poolConfig) ConnectTimeout() time.Duration  { return p.Raw.ConnectTimeout }
func (p *poolConfig) ShutdownTimeout() time.Duration { return p.Raw.ShutdownTimeout }
func (p *poolConfig) MaxConnecting() int             { return p.Raw.MaxConnecting }
