package services

import (
	"fmt"
	"time"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

var _ contracts.ServicesConfig = (servicesConfig)(nil)

// servicesConfig — map именованных сервисов.
type servicesConfig map[string]*serviceConfig

// fullConfig нужен, чтобы Viper мог распарсить секцию services.
type fullConfig struct {
	Services map[string]rawService `mapstructure:"services" yaml:"services"`
}

// NewServicesConfig создает конфигурацию внешних сервисов. YAML -> ENV.
func newServicesConfig() (contracts.ServicesConfig, error) {
	if v := helpers.GetViper(); v != nil {
		var full fullConfig
		if err := v.Unmarshal(&full); err == nil && len(full.Services) > 0 {
			cfg := make(servicesConfig)
			for name, raw := range full.Services {
				if raw.Timeout == 0 {
					raw.Timeout = 30 * time.Second
				}

				svc, err := newServiceConfig(raw)
				if err != nil {
					return nil, fmt.Errorf("service %s: %w", name, err)
				}
				cfg[name] = svc
			}
			return cfg, nil
		}
	}

	return servicesConfig{}, nil
}

// Get возвращает сервис по имени.
func (sc servicesConfig) Get(name string) (contracts.ServiceConfig, bool) {
	svc, exists := sc[name]
	return svc, exists
}

// All возвращает все сервисы как map.
func (sc servicesConfig) All() map[string]contracts.ServiceConfig {
	out := make(map[string]contracts.ServiceConfig, len(sc))
	for name, svc := range sc {
		out[name] = svc
	}
	return out
}
