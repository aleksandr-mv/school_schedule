package telegram

import (
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.TelegramConfig = (*Config)(nil)

// rawConfig для загрузки данных из YAML/ENV
type rawConfig struct {
	Token      string `mapstructure:"token"         yaml:"token"         env:"TELEGRAM_TOKEN"`
	WebhookURL string `mapstructure:"webhook_url"   yaml:"webhook_url"   env:"TELEGRAM_WEBHOOK_URL"`
	DebugMode  bool   `mapstructure:"debug_mode"    yaml:"debug_mode"    env:"TELEGRAM_DEBUG_MODE"`
	Timeout    int    `mapstructure:"timeout"       yaml:"timeout"       env:"TELEGRAM_TIMEOUT"`
	MaxRetries int    `mapstructure:"max_retries"   yaml:"max_retries"   env:"TELEGRAM_MAX_RETRIES"`
}

// Config публичная структура Telegram конфигурации
type Config struct {
	raw rawConfig
}

// defaultConfig возвращает rawConfig с дефолтными значениями
func defaultConfig() rawConfig {
	return rawConfig{
		Token:      "",
		WebhookURL: "",
		DebugMode:  false,
		Timeout:    30,
		MaxRetries: 3,
	}
}

// Методы для TelegramConfig интерфейса
func (c *Config) IsEnabled() bool    { return c.raw.Token != "" }
func (c *Config) Token() string      { return c.raw.Token }
func (c *Config) WebhookURL() string { return c.raw.WebhookURL }
func (c *Config) DebugMode() bool    { return c.raw.DebugMode }
func (c *Config) Timeout() int       { return c.raw.Timeout }
func (c *Config) MaxRetries() int    { return c.raw.MaxRetries }
