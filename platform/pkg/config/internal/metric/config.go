package metric

import (
	"time"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.MetricConfig = (*Config)(nil)

// rawConfig для загрузки данных из YAML/ENV
type rawConfig struct {
	Enable          bool          `mapstructure:"enable" yaml:"enable" env:"METRIC_ENABLE"`
	Endpoint        string        `mapstructure:"endpoint" yaml:"endpoint" env:"METRIC_ENDPOINT"`
	Timeout         time.Duration `mapstructure:"timeout" yaml:"timeout" env:"METRIC_TIMEOUT"`
	Namespace       string        `mapstructure:"namespace" yaml:"namespace" env:"METRIC_NAMESPACE"`
	AppName         string        `mapstructure:"app_name" yaml:"app_name" env:"METRIC_APP_NAME"`
	ExportInterval  time.Duration `mapstructure:"export_interval" yaml:"export_interval" env:"METRIC_EXPORT_INTERVAL"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout" env:"METRIC_SHUTDOWN_TIMEOUT"`

	// Bucket'ы для гистограмм
	BucketBoundaries []float64 `mapstructure:"bucket_boundaries" yaml:"bucket_boundaries"`
}

// Config публичная структура Metric конфигурации
type Config struct {
	raw rawConfig
}

// defaultConfig возвращает rawConfig с дефолтными значениями
func defaultConfig() rawConfig {
	return rawConfig{
		Enable:          false,
		Endpoint:        "localhost:4317",
		Timeout:         30 * time.Second,
		Namespace:       "microservices",
		AppName:         "platform",
		ExportInterval:  5 * time.Second,
		ShutdownTimeout: 30 * time.Second,

		// Bucket'ы по умолчанию для гистограмм (от 0.1ms до ~3.3s)
		BucketBoundaries: []float64{
			0.0001, 0.0002, 0.0004, 0.0008, 0.0016, 0.0032, 0.0064, 0.0128,
			0.0256, 0.0512, 0.1024, 0.2048, 0.4096, 0.8192, 1.6384, 3.2768,
		},
	}
}

// Методы для MetricConfig интерфейса
func (c *Config) Enable() bool                   { return c.raw.Enable }
func (c *Config) Endpoint() string               { return c.raw.Endpoint }
func (c *Config) Timeout() time.Duration         { return c.raw.Timeout }
func (c *Config) Namespace() string              { return c.raw.Namespace }
func (c *Config) AppName() string                { return c.raw.AppName }
func (c *Config) ExportInterval() time.Duration  { return c.raw.ExportInterval }
func (c *Config) ShutdownTimeout() time.Duration { return c.raw.ShutdownTimeout }
func (c *Config) BucketBoundaries() []float64    { return c.raw.BucketBoundaries }
