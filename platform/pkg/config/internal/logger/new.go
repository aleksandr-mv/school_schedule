package logger

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/helpers"
)

// New создает Logger конфигурацию по стратегии: Defaults → YAML → ENV
func New() (contracts.LoggerConfig, error) {
	// 1. Создаем конфигурацию с дефолтными значениями
	cfg := &Config{
		raw: defaultConfig(),
	}

	// 2. Перезаписываем YAML'ом (если есть)
	if section := helpers.GetSection("logger"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("failed to unmarshal logger YAML: %w", err)
		}
	}

	// 3. Перезаписываем ENV переменными (финальный приоритет)
	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("failed to parse logger ENV: %w", err)
	}

	return cfg, nil
}
