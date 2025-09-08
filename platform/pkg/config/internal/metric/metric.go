package metric

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// rawMetricConfig содержит настройки метрик.
type rawMetricConfig struct {
	EnableFlag   bool          `mapstructure:"enable" yaml:"enable" env:"METRIC_ENABLE"`
	EndpointAddr string        `mapstructure:"endpoint" yaml:"endpoint" env:"METRIC_ENDPOINT"`
	TimeoutDur   time.Duration `mapstructure:"timeout" yaml:"timeout" env:"METRIC_TIMEOUT"`

	NamespaceStr string `mapstructure:"namespace" yaml:"namespace" env:"METRIC_NAMESPACE"`
	AppNameStr   string `mapstructure:"app_name" yaml:"app_name" env:"METRIC_APP_NAME"`

	ExportIntervalDur  time.Duration `mapstructure:"export_interval" yaml:"export_interval" env:"METRIC_EXPORT_INTERVAL"`
	ShutdownTimeoutDur time.Duration `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout" env:"METRIC_SHUTDOWN_TIMEOUT"`

	// Bucket'ы для гистограмм
	BucketBoundariesSlice []float64 `mapstructure:"bucket_boundaries" yaml:"bucket_boundaries"`
}

// defaultMetricConfig возвращает конфигурацию с дефолтными значениями
func defaultMetricConfig() *rawMetricConfig {
	return &rawMetricConfig{
		EnableFlag:         false,
		EndpointAddr:       "localhost:4317",
		TimeoutDur:         30 * time.Second,
		NamespaceStr:       "microservices",
		AppNameStr:         "platform",
		ExportIntervalDur:  5 * time.Second,
		ShutdownTimeoutDur: 30 * time.Second,

		// Bucket'ы по умолчанию для гистограмм (от 0.1ms до ~3.3s)
		BucketBoundariesSlice: []float64{
			0.0001, 0.0002, 0.0004, 0.0008, 0.0016, 0.0032, 0.0064, 0.0128,
			0.0256, 0.0512, 0.1024, 0.2048, 0.4096, 0.8192, 1.6384, 3.2768,
		},
	}
}

// DefaultMetricConfig читает конфигурацию метрик из ENV.
func DefaultMetricConfig() (*rawMetricConfig, error) {
	raw := defaultMetricConfig()
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse metric env: %w", err)
	}

	return raw, nil
}

// NewMetricConfig создает конфигурацию метрик, пытаясь сначала загрузить из YAML, затем из ENV.
func NewMetricConfig() (*rawMetricConfig, error) {
	if section := helpers.GetSection("metric"); section != nil {
		raw := defaultMetricConfig()
		if err := section.Unmarshal(raw); err == nil {
			return raw, nil
		}
	}

	return DefaultMetricConfig()
}
