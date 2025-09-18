package mongo

import (
	"time"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.MongoPoolConfig = (*Pool)(nil)

// rawPool для загрузки данных из YAML/ENV
type rawPool struct {
	MaxSize         int           `mapstructure:"max_pool_size"      yaml:"max_pool_size"      env:"MONGO_MAX_POOL_SIZE"`
	MinSize         int           `mapstructure:"min_pool_size"      yaml:"min_pool_size"      env:"MONGO_MIN_POOL_SIZE"`
	MaxConnecting   int           `mapstructure:"max_connecting"     yaml:"max_connecting"     env:"MONGO_MAX_CONNECTING"`
	MaxIdleTime     time.Duration `mapstructure:"max_conn_idle_time" yaml:"max_conn_idle_time" env:"MONGO_MAX_CONN_IDLE_TIME"`
	ConnectTimeout  time.Duration `mapstructure:"connect_timeout"    yaml:"connect_timeout"    env:"MONGO_CONNECT_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"   yaml:"shutdown_timeout"   env:"MONGO_SHUTDOWN_TIMEOUT"`
	HeartbeatPeriod time.Duration `mapstructure:"heartbeat_period"   yaml:"heartbeat_period"   env:"MONGO_HEARTBEAT_PERIOD"`
	SelectTimeout   time.Duration `mapstructure:"server_select_timeout" yaml:"server_select_timeout" env:"MONGO_SERVER_SELECT_TIMEOUT"`
}

// Pool публичная структура для использования
type Pool struct {
	raw rawPool
}

// defaultPool возвращает rawPool с дефолтными значениями
func defaultPool() rawPool {
	return rawPool{
		MaxSize:         10,
		MinSize:         2,
		MaxConnecting:   5,
		MaxIdleTime:     30 * time.Minute,
		ConnectTimeout:  5 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		HeartbeatPeriod: 10 * time.Second,
		SelectTimeout:   30 * time.Second,
	}
}

// Методы для PoolConfig интерфейса (базовые)
func (p *Pool) MaxSize() int                   { return p.raw.MaxSize }
func (p *Pool) MinSize() int                   { return p.raw.MinSize }
func (p *Pool) MaxIdleTime() time.Duration     { return p.raw.MaxIdleTime }
func (p *Pool) ConnectTimeout() time.Duration  { return p.raw.ConnectTimeout }
func (p *Pool) ShutdownTimeout() time.Duration { return p.raw.ShutdownTimeout }

// Методы для MongoPoolConfig интерфейса (специфичные для MongoDB)
func (p *Pool) MaxConnecting() int                    { return p.raw.MaxConnecting }
func (p *Pool) HeartbeatPeriod() time.Duration        { return p.raw.HeartbeatPeriod }
func (p *Pool) ServerSelectionTimeout() time.Duration { return p.raw.SelectTimeout }
