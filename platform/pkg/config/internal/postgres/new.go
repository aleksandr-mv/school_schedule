package postgres

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// New создает PostgreSQL конфигурацию по стратегии: Defaults → YAML → ENV
func New() (contracts.PostgresConfig, error) {
	// 1. Создаем конфигурацию с дефолтными значениями
	cfg := &Config{
		raw: defaultConfig(),
	}

	// 2. Перезаписываем YAML'ом (если есть)
	if section := helpers.GetSection("postgres"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("failed to unmarshal postgres YAML: %w", err)
		}
	}

	// 3. Перезаписываем ENV переменными (финальный приоритет)
	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("failed to parse postgres ENV: %w", err)
	}

	// 4. Инициализируем подструктуры
	cfg.primaryConn = &Connection{raw: cfg.raw.Primary}
	cfg.replicaConns = make([]*Connection, len(cfg.raw.Replicas))
	for i, replica := range cfg.raw.Replicas {
		cfg.replicaConns[i] = &Connection{raw: replica}
	}
	cfg.poolSettings = &Pool{raw: cfg.raw.Pool}

	return cfg, nil
}
