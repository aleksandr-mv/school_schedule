package builder

import (
	redigo "github.com/gomodule/redigo/redis"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// RedisBuilder создает Redis клиенты из конфигурации
type RedisBuilder struct {
	config contracts.RedisConfig
}

// NewRedisBuilder создает билдер для Redis
func NewRedisBuilder(config contracts.RedisConfig) *RedisBuilder {
	return &RedisBuilder{
		config: config,
	}
}

// BuildPool создает Redis пул
func (b *RedisBuilder) BuildPool() (*redigo.Pool, error) {
	if !b.config.IsEnabled() {
		return nil, nil
	}

	connection := b.config.Connection()
	pool := b.config.Pool()

	redigoPool := &redigo.Pool{
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", connection.Address(),
				redigo.DialPassword(connection.Password()),
				redigo.DialDatabase(connection.Database()),
				redigo.DialConnectTimeout(pool.DialTimeout()),
				redigo.DialReadTimeout(pool.ReadTimeout()),
				redigo.DialWriteTimeout(pool.WriteTimeout()),
			)
		},
		MaxIdle:     pool.MinIdle(),
		MaxActive:   pool.PoolSize(),
		IdleTimeout: pool.IdleTimeout(),
		Wait:        true,
	}

	return redigoPool, nil
}
