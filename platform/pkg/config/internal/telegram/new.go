package telegram

import (
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

type module struct {
	config *rawTelegramConfig
}

// New создает модуль конфигурации Telegram.
func New() (contracts.TelegramConfig, error) {
	telegramCfg, err := NewTelegramConfig()
	if err != nil {
		return nil, err
	}

	return &module{config: telegramCfg}, nil
}

func (m *module) IsEnabled() bool {
	return m.config.Token != ""
}

func (m *module) Token() string {
	return m.config.Token
}

func (m *module) WebhookURL() string {
	return m.config.WebhookURL
}

func (m *module) DebugMode() bool {
	return m.config.DebugMode
}

func (m *module) Timeout() int {
	return m.config.Timeout
}

func (m *module) MaxRetries() int {
	return m.config.MaxRetries
}
