package telegram

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// rawTelegramConfig соответствует полям секции telegram в YAML и env-переменным.
type rawTelegramConfig struct {
	Token      string `mapstructure:"token"         yaml:"token"         env:"TELEGRAM_TOKEN"`
	WebhookURL string `mapstructure:"webhook_url"   yaml:"webhook_url"   env:"TELEGRAM_WEBHOOK_URL"`
	DebugMode  bool   `mapstructure:"debug_mode"    yaml:"debug_mode"    env:"TELEGRAM_DEBUG_MODE"`
	Timeout    int    `mapstructure:"timeout"       yaml:"timeout"       env:"TELEGRAM_TIMEOUT"`
	MaxRetries int    `mapstructure:"max_retries"   yaml:"max_retries"   env:"TELEGRAM_MAX_RETRIES"`
}

// defaultTelegramConfig возвращает конфигурацию Telegram с дефолтными значениями
func defaultTelegramConfig() *rawTelegramConfig {
	return &rawTelegramConfig{
		Token:      "",
		WebhookURL: "",
		DebugMode:  false,
		Timeout:    30,
		MaxRetries: 3,
	}
}

// DefaultTelegramConfig читает конфигурацию Telegram из ENV.
func DefaultTelegramConfig() (*rawTelegramConfig, error) {
	// Начинаем с дефолтной конфигурации
	raw := defaultTelegramConfig()

	// Применяем переменные окружения поверх дефолтов
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse telegram env: %w", err)
	}

	return raw, nil
}

// NewTelegramConfig создает конфигурацию Telegram, пытаясь сначала загрузить из YAML, затем из ENV.
func NewTelegramConfig() (*rawTelegramConfig, error) {
	if section := helpers.GetSection("telegram"); section != nil {
		// Начинаем с дефолтной конфигурации
		raw := defaultTelegramConfig()

		// Применяем YAML конфигурацию поверх дефолтов
		if err := section.Unmarshal(raw); err == nil {
			return raw, nil
		}
	}

	return DefaultTelegramConfig()
}
