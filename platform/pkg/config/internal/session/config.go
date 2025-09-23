package session

import (
	"time"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.SessionConfig = (*Config)(nil)

// rawConfig для загрузки данных из YAML/ENV
type rawConfig struct {
	TTL time.Duration `mapstructure:"ttl" yaml:"ttl" env:"SESSION_TTL"`
}

// Config публичная структура Session конфигурации
type Config struct {
	raw rawConfig
}

// defaultConfig возвращает rawConfig с дефолтными значениями
func defaultConfig() rawConfig {
	return rawConfig{
		TTL: 24 * time.Hour,
	}
}

// Методы для SessionConfig интерфейса
func (c *Config) TTL() time.Duration {
	return c.raw.TTL
}
