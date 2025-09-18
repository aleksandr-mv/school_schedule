package tracing

import (
	"time"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.TracingConfig = (*Config)(nil)

// rawConfig для загрузки данных из YAML/ENV
type rawConfig struct {
	Enable   bool          `mapstructure:"enable" yaml:"enable" env:"TRACING_ENABLE"`
	Endpoint string        `mapstructure:"endpoint" yaml:"endpoint" env:"TRACING_ENDPOINT"`
	Timeout  time.Duration `mapstructure:"timeout" yaml:"timeout" env:"TRACING_TIMEOUT"`

	SampleRatio int `mapstructure:"sample_ratio" yaml:"sample_ratio" env:"TRACING_SAMPLE_RATIO"`

	RetryEnabled         bool          `mapstructure:"retry_enabled" yaml:"retry_enabled" env:"TRACING_RETRY_ENABLED"`
	RetryInitialInterval time.Duration `mapstructure:"retry_initial_interval" yaml:"retry_initial_interval" env:"TRACING_RETRY_INITIAL_INTERVAL"`
	RetryMaxInterval     time.Duration `mapstructure:"retry_max_interval" yaml:"retry_max_interval" env:"TRACING_RETRY_MAX_INTERVAL"`
	RetryMaxElapsedTime  time.Duration `mapstructure:"retry_max_elapsed_time" yaml:"retry_max_elapsed_time" env:"TRACING_RETRY_MAX_ELAPSED_TIME"`

	EnableTraceContext bool          `mapstructure:"enable_trace_context" yaml:"enable_trace_context" env:"TRACING_ENABLE_TRACE_CONTEXT"`
	EnableBaggage      bool          `mapstructure:"enable_baggage" yaml:"enable_baggage" env:"TRACING_ENABLE_BAGGAGE"`
	ShutdownTimeout    time.Duration `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout" env:"TRACING_SHUTDOWN_TIMEOUT"`
}

// Config публичная структура Tracing конфигурации
type Config struct {
	raw rawConfig
}

// defaultConfig возвращает rawConfig с дефолтными значениями
func defaultConfig() rawConfig {
	return rawConfig{
		Enable:               true,
		Endpoint:             "localhost:4317",
		Timeout:              5 * time.Second,
		SampleRatio:          100,
		RetryEnabled:         true,
		RetryInitialInterval: 500 * time.Millisecond,
		RetryMaxInterval:     5 * time.Second,
		RetryMaxElapsedTime:  30 * time.Second,
		EnableTraceContext:   true,
		EnableBaggage:        true,
		ShutdownTimeout:      30 * time.Second,
	}
}

// Методы для TracingConfig интерфейса
func (c *Config) Enable() bool                        { return c.raw.Enable }
func (c *Config) Endpoint() string                    { return c.raw.Endpoint }
func (c *Config) Timeout() time.Duration              { return c.raw.Timeout }
func (c *Config) SampleRatio() int                    { return c.raw.SampleRatio }
func (c *Config) RetryEnabled() bool                  { return c.raw.RetryEnabled }
func (c *Config) RetryInitialInterval() time.Duration { return c.raw.RetryInitialInterval }
func (c *Config) RetryMaxInterval() time.Duration     { return c.raw.RetryMaxInterval }
func (c *Config) RetryMaxElapsedTime() time.Duration  { return c.raw.RetryMaxElapsedTime }
func (c *Config) EnableTraceContext() bool            { return c.raw.EnableTraceContext }
func (c *Config) EnableBaggage() bool                 { return c.raw.EnableBaggage }
func (c *Config) ShutdownTimeout() time.Duration      { return c.raw.ShutdownTimeout }
