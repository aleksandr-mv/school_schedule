package redis

import "github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"

type module struct {
	redis contracts.RedisConfig
}

// New создает модуль конфигурации Redis.
func New() (contracts.RedisConfig, error) {
	redisCfg, err := NewRedisConfig()
	if err != nil {
		return nil, err
	}

	return &module{redis: redisCfg}, nil
}

func (m *module) IsEnabled() bool {
	return m.redis.IsEnabled()
}

func (m *module) Connection() contracts.RedisConnection {
	return m.redis.Connection()
}

func (m *module) Pool() contracts.RedisPoolConfig {
	return m.redis.Pool()
}
