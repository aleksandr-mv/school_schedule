package redis

import (
	"time"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.RedisPoolConfig = (*Pool)(nil)

// rawPool для загрузки данных из YAML/ENV
type rawPool struct {
	MaxActive    int           `mapstructure:"max_active"      yaml:"max_active"      env:"REDIS_POOL_MAX_ACTIVE"`
	MaxIdle      int           `mapstructure:"max_idle"        yaml:"max_idle"        env:"REDIS_POOL_MAX_IDLE"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"    yaml:"idle_timeout"    env:"REDIS_POOL_IDLE_TIMEOUT"`
	ConnTimeout  time.Duration `mapstructure:"conn_timeout"    yaml:"conn_timeout"    env:"REDIS_POOL_CONN_TIMEOUT"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"    yaml:"read_timeout"    env:"REDIS_POOL_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"   yaml:"write_timeout"   env:"REDIS_POOL_WRITE_TIMEOUT"`
	PoolTimeout  time.Duration `mapstructure:"pool_timeout"    yaml:"pool_timeout"    env:"REDIS_POOL_TIMEOUT"`
}

// Pool публичная структура для использования
type Pool struct {
	raw rawPool
}

// defaultPool возвращает rawPool с дефолтными значениями
func defaultPool() rawPool {
	return rawPool{
		MaxActive:    10,
		MaxIdle:      5,
		IdleTimeout:  240 * time.Second,
		ConnTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	}
}

// Методы для RedisPoolConfig интерфейса
func (p *Pool) MaxActive() int              { return p.raw.MaxActive }
func (p *Pool) MaxIdle() int                { return p.raw.MaxIdle }
func (p *Pool) IdleTimeout() time.Duration  { return p.raw.IdleTimeout }
func (p *Pool) ConnTimeout() time.Duration  { return p.raw.ConnTimeout }
func (p *Pool) ReadTimeout() time.Duration  { return p.raw.ReadTimeout }
func (p *Pool) WriteTimeout() time.Duration { return p.raw.WriteTimeout }
func (p *Pool) PoolTimeout() time.Duration  { return p.raw.PoolTimeout }
