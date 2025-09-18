package services

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// New создает Services конфигурацию по стратегии: Defaults → YAML → ENV
func New() (contracts.ServicesConfig, error) {
	// 1. Создаем конфигурацию с дефолтными значениями
	cfg := &Config{
		raw:      defaultConfig(),
		services: make(map[string]*ServiceConfig),
	}

	// 2. Перезаписываем YAML'ом (если есть)
	if section := helpers.GetSection("services"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("failed to unmarshal services YAML: %w", err)
		}
	}

	// 3. Перезаписываем ENV переменными (финальный приоритет)
	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("failed to parse services ENV: %w", err)
	}

	// 4. Инициализируем сервисы из raw данных
	for name, rawSvc := range cfg.raw.Services {
		// Применяем дефолтный таймаут если не задан
		if rawSvc.Timeout == 0 {
			rawSvc.Timeout = defaultService().Timeout
		}
		cfg.services[name] = newServiceConfig(rawSvc)
	}

	return cfg, nil
}
