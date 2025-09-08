package logger

import "github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"

// rawOTLPConfig содержит OTLP настройки для интеграции с OpenTelemetry.
type rawOTLPConfig struct {
	Enable          bool   `mapstructure:"enable" yaml:"enable" env:"ENABLE"`
	Endpoint        string `mapstructure:"endpoint" yaml:"endpoint" env:"ENDPOINT"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT"`
}

// defaultOTLPConfig возвращает OTLP конфигурацию с дефолтными значениями
func defaultOTLPConfig() *rawOTLPConfig {
	return &rawOTLPConfig{
		Enable:          false,
		Endpoint:        "localhost:4317",
		ShutdownTimeout: 2,
	}
}

// otlpModule реализует OTLPModule интерфейс.
type otlpModule struct {
	config *rawOTLPConfig
}

// NewOTLPModule создает новый OTLP модуль.
func NewOTLPModule(config *rawOTLPConfig) contracts.OTLPModule {
	return &otlpModule{config: config}
}

// Методы OTLPModule
func (m *otlpModule) Enable() bool         { return m.config.Enable }
func (m *otlpModule) Endpoint() string     { return m.config.Endpoint }
func (m *otlpModule) ShutdownTimeout() int { return m.config.ShutdownTimeout }
