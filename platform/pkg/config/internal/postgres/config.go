package postgres

import (
	"math/rand"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.PostgresConfig = (*Config)(nil)

// rawConfig для загрузки данных из YAML/ENV
type rawConfig struct {
	Primary  rawConnection   `mapstructure:"primary"  yaml:"primary"`
	Replicas []rawConnection `mapstructure:"replicas" yaml:"replicas"`
	Pool     rawPool         `mapstructure:"pool"     yaml:"pool"`
}

// Config публичная структура PostgreSQL конфигурации с Primary-Replica
type Config struct {
	raw          rawConfig
	primaryConn  *Connection
	replicaConns []*Connection
	poolSettings *Pool
}

// defaultConfig возвращает rawConfig с дефолтными значениями
func defaultConfig() rawConfig {
	return rawConfig{
		Primary:  defaultConnection(),
		Replicas: []rawConnection{defaultConnection()}, // По умолчанию реплика = primary
		Pool:     defaultPool(),
	}
}

// Primary возвращает primary соединение
func (c *Config) Primary() contracts.DBConnection { return c.primaryConn }

// PrimaryURI возвращает URI для записи
func (c *Config) PrimaryURI() string { return c.primaryConn.URI() }

// Replicas возвращает все реплики
func (c *Config) Replicas() []contracts.DBConnection {
	replicas := make([]contracts.DBConnection, len(c.replicaConns))
	for i := range c.replicaConns {
		replicas[i] = c.replicaConns[i]
	}
	return replicas
}

// ReplicaURI возвращает случайную реплику
func (c *Config) ReplicaURI() string {
	if len(c.replicaConns) == 0 {
		return c.PrimaryURI() // fallback
	}
	replica := c.replicaConns[rand.Intn(len(c.replicaConns))]
	return replica.URI()
}

// Pool возвращает настройки пула
func (c *Config) Pool() contracts.PostgresPoolConfig { return c.poolSettings }
