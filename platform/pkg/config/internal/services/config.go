package services

import (
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.ServicesConfig = (*Config)(nil)

// Config публичная структура Services конфигурации
type Config struct {
	services map[string]*ServiceConfig
}

// Get возвращает сервис по имени
func (c *Config) Get(name string) (contracts.ServiceConfig, bool) {
	svc, exists := c.services[name]
	return svc, exists
}

// All возвращает все сервисы как map
func (c *Config) All() map[string]contracts.ServiceConfig {
	out := make(map[string]contracts.ServiceConfig, len(c.services))
	for name, svc := range c.services {
		out[name] = svc
	}
	return out
}
