package builder

import (
	"github.com/redis/go-redis/v9"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/cache"
	redisclient "github.com/aleksandr-mv/school_schedule/platform/pkg/cache/redis"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

type RedisBuilder struct {
	config contracts.RedisConfig
}

func NewRedisBuilder(config contracts.RedisConfig) *RedisBuilder {
	return &RedisBuilder{
		config: config,
	}
}

// BuildClient создает Redis кластер клиент через go-redis
func (b *RedisBuilder) BuildClient() (cache.RedisClient, error) {
	cluster := b.config.Cluster()
	if !cluster.IsEnabled() {
		return nil, nil
	}

	pool := b.config.Pool()

	// Создаем кластерный клиент
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          cluster.Nodes(),
		Password:       cluster.Password(),
		MaxRedirects:   cluster.MaxRedirects(),
		ReadOnly:       cluster.ReadOnlyCommands(),
		RouteByLatency: cluster.RouteByLatency(),
		RouteRandomly:  cluster.RouteRandomly(),

		// Настройки пула
		PoolSize:        pool.MaxActive(),
		MinIdleConns:    pool.MaxIdle(),
		PoolTimeout:     pool.PoolTimeout(),
		ConnMaxIdleTime: pool.IdleTimeout(),
		DialTimeout:     pool.ConnTimeout(),
		ReadTimeout:     pool.ReadTimeout(),
		WriteTimeout:    pool.WriteTimeout(),
	})

	// Создаем адаптер с нашим интерфейсом
	return redisclient.NewClient(rdb, logger.Logger(), pool.ConnTimeout()), nil
}
