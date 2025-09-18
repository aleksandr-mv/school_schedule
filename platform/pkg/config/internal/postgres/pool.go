package postgres

import (
	"time"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.PostgresPoolConfig = (*Pool)(nil)

// rawPool для загрузки данных из YAML/ENV
type rawPool struct {
	MaxConns     int           `mapstructure:"max_cons"            yaml:"max_cons"            env:"DB_MAX_CONS"`
	MinConns     int           `mapstructure:"min_cons"            yaml:"min_cons"            env:"DB_MIN_CONS"`
	MaxLifetime  time.Duration `mapstructure:"max_conn_lifetime"   yaml:"max_conn_lifetime"   env:"DB_MAX_CONN_LIFETIME"`
	MaxIdle      time.Duration `mapstructure:"max_conn_idle_time"  yaml:"max_conn_idle_time"  env:"DB_MAX_CONN_IDLE_TIME"`
	HealthPeriod time.Duration `mapstructure:"health_check_period" yaml:"health_check_period" env:"DB_HEALTH_CHECK_PERIOD"`
	ConnTimeout  time.Duration `mapstructure:"connect_timeout"     yaml:"connect_timeout"     env:"DB_CONNECT_TIMEOUT"`
	ShutTimeout  time.Duration `mapstructure:"shutdown_timeout"    yaml:"shutdown_timeout"    env:"DB_SHUTDOWN_TIMEOUT"`
}

// Pool публичная структура для использования
type Pool struct {
	raw rawPool
}

// defaultPool возвращает rawPool с дефолтными значениями
func defaultPool() rawPool {
	return rawPool{
		MaxConns:     10,
		MinConns:     2,
		MaxLifetime:  time.Hour,
		MaxIdle:      30 * time.Minute,
		HealthPeriod: time.Minute,
		ConnTimeout:  5 * time.Second,
		ShutTimeout:  30 * time.Second,
	}
}

// Методы для PostgresPoolConfig интерфейса
func (p *Pool) MaxSize() int                     { return p.raw.MaxConns }
func (p *Pool) MinSize() int                     { return p.raw.MinConns }
func (p *Pool) MaxIdleTime() time.Duration       { return p.raw.MaxIdle }
func (p *Pool) ConnectTimeout() time.Duration    { return p.raw.ConnTimeout }
func (p *Pool) ShutdownTimeout() time.Duration   { return p.raw.ShutTimeout }
func (p *Pool) MaxConnLifetime() time.Duration   { return p.raw.MaxLifetime }
func (p *Pool) HealthCheckPeriod() time.Duration { return p.raw.HealthPeriod }
