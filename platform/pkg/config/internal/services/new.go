package services

import (
	"fmt"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// New создает Services конфигурацию из YAML
func New() (contracts.ServicesConfig, error) {
	// 1. Создаем конфигурацию с пустой картой сервисов
	cfg := &Config{
		services: make(map[string]*ServiceConfig),
	}

	// 2. Загружаем сервисы из YAML (если есть)
	if section := helpers.GetSection("services"); section != nil {
		var servicesMap map[string]rawServiceConfig
		if err := section.Unmarshal(&servicesMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal services YAML: %w", err)
		}

		// 3. Инициализируем сервисы из YAML данных
		for name, rawSvc := range servicesMap {
			// Применяем дефолтный таймаут если не задан
			if rawSvc.Timeout == 0 {
				rawSvc.Timeout = defaultService().Timeout
			}
			cfg.services[name] = newServiceConfig(rawSvc)
		}
	}

	return cfg, nil
}
