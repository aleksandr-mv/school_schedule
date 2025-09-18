package redis

import (
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.RedisConfig = (*Config)(nil)

// rawConfig для загрузки данных из YAML/ENV
type rawConfig struct {
	Cluster rawCluster `mapstructure:"cluster" yaml:"cluster"`
	Pool    rawPool    `mapstructure:"pool"    yaml:"pool"`
}

// Config публичная структура Redis конфигурации (только кластер)
type Config struct {
	raw           rawConfig
	clusterConfig *Cluster
	poolConfig    *Pool
}

// defaultConfig возвращает rawConfig с дефолтными значениями
func defaultConfig() rawConfig {
	return rawConfig{
		Cluster: defaultCluster(),
		Pool:    defaultPool(),
	}
}

// Методы для RedisConfig интерфейса
func (c *Config) Cluster() contracts.RedisClusterConfig {
	if c.clusterConfig == nil {
		c.clusterConfig = &Cluster{raw: c.raw.Cluster}
	}
	return c.clusterConfig
}

func (c *Config) Pool() contracts.RedisPoolConfig {
	if c.poolConfig == nil {
		c.poolConfig = &Pool{raw: c.raw.Pool}
	}
	return c.poolConfig
}
