package tracing

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// rawTracingConfig содержит настройки трейсинга.
type rawTracingConfig struct {
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

// defaultTracingConfig возвращает конфигурацию с дефолтными значениями
func defaultTracingConfig() *rawTracingConfig {
	return &rawTracingConfig{
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

// DefaultTracingConfig читает конфигурацию трассировки из ENV.
func DefaultTracingConfig() (*rawTracingConfig, error) {
	// Начинаем с дефолтной конфигурации
	raw := defaultTracingConfig()

	// Применяем переменные окружения поверх дефолтов
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse tracing env: %w", err)
	}

	return raw, nil
}

// NewTracingConfig создает конфигурацию трассировки, пытаясь сначала загрузить из YAML, затем из ENV.
func NewTracingConfig() (*rawTracingConfig, error) {
	if section := helpers.GetSection("tracing"); section != nil {
		// Начинаем с дефолтной конфигурации
		raw := defaultTracingConfig()

		// Применяем YAML конфигурацию поверх дефолтов
		if err := section.Unmarshal(raw); err == nil {
			return raw, nil
		}
	}

	return DefaultTracingConfig()
}
