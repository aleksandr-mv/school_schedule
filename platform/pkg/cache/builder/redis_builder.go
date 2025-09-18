package builder

import (
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/go-redis/redis/v8"
)

type RedisBuilder struct {
	config contracts.RedisConfig
}

func NewRedisBuilder(config contracts.RedisConfig) *RedisBuilder {
	return &RedisBuilder{
		config: config,
	}
}

// BuildClient создает Redis кластер клиент (go-redis)
func (b *RedisBuilder) BuildClient() (redis.UniversalClient, error) {
	cluster := b.config.Cluster()
	if !cluster.IsEnabled() {
		return nil, nil
	}

	pool := b.config.Pool()

	// Создаем кластерный клиент
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          cluster.Nodes(),
		Password:       cluster.Password(),
		MaxRedirects:   cluster.MaxRedirects(),
		ReadOnly:       cluster.ReadOnlyCommands(),
		RouteByLatency: cluster.RouteByLatency(),
		RouteRandomly:  cluster.RouteRandomly(),

		// Настройки пула
		PoolSize:     pool.MaxActive(),
		MinIdleConns: pool.MaxIdle(),
		IdleTimeout:  pool.IdleTimeout(),
		DialTimeout:  pool.ConnTimeout(),
		ReadTimeout:  pool.ReadTimeout(),
		WriteTimeout: pool.WriteTimeout(),
	})

	return client, nil
}
