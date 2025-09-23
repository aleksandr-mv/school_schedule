package app

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/helpers"
)

// New создает App конфигурацию по стратегии: Defaults → YAML → ENV
func New() (contracts.AppConfig, error) {
	// 1. Создаем конфигурацию с дефолтными значениями
	cfg := &Config{
		raw: defaultConfig(),
	}

	// 2. Перезаписываем YAML'ом (если есть)
	if section := helpers.GetSection("app"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("failed to unmarshal app YAML: %w", err)
		}
	}

	// 3. Перезаписываем ENV переменными (финальный приоритет)
	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("failed to parse app ENV: %w", err)
	}

	return cfg, nil
}
