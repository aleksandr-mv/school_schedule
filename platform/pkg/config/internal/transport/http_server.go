package transport

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// rawHTTPConfig соответствует полям внутри секции http_server.
type rawHTTPConfig struct {
	Host              string        `mapstructure:"host"             yaml:"host"               env:"HTTP_HOST"`
	Port              int           `mapstructure:"port"             yaml:"port"               env:"HTTP_PORT"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout" yaml:"read_header_timeout" env:"HTTP_READ_TIMEOUT"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"     yaml:"read_timeout"       env:"HTTP_READ_TIMEOUT"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"    yaml:"write_timeout"      env:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"     yaml:"idle_timeout"       env:"HTTP_IDLE_TIMEOUT"`
	MaxHeaderBytes    int           `mapstructure:"max_header_bytes" yaml:"max_header_bytes"   env:"HTTP_MAX_HEADER_BYTES"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout"   env:"HTTP_SHUTDOWN_TIMEOUT"`
	HandlerTimeout    time.Duration `mapstructure:"handler_timeout"  yaml:"handler_timeout"    env:"HTTP_HANDLER_TIMEOUT"`
}

// httpServerConfig хранит данные секции http_server и реализует HTTPConfig.
// Тег yaml:"http_server" применяется сразу к полю raw.
type httpServerConfig struct {
	Raw rawHTTPConfig `yaml:"http_server"`
}

var _ contracts.HTTPServer = (*httpServerConfig)(nil)

// defaultHTTPServerConfig возвращает конфигурацию HTTP server с дефолтными значениями
func defaultHTTPServerConfig() *rawHTTPConfig {
	return &rawHTTPConfig{
		Host:              "localhost",
		Port:              8080,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
		ShutdownTimeout:   10 * time.Second,
		HandlerTimeout:    0, // будет использовать ReadTimeout
	}
}

// DefaultHTTPServerConfig читает конфиг HTTP из ENV
func DefaultHTTPServerConfig() (*httpServerConfig, error) {
	// Начинаем с дефолтной конфигурации
	raw := defaultHTTPServerConfig()

	// Применяем переменные окружения поверх дефолтов
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse HTTP server env: %w", err)
	}

	return &httpServerConfig{Raw: *raw}, nil
}

// NewHTTPServerConfig создает конфигурацию HTTP сервера. YAML -> ENV -> валидация.
func NewHTTPServerConfig() (*httpServerConfig, error) {
	// Попытка загрузки из YAML
	if section := helpers.GetSection("http_server"); section != nil {
		// Начинаем с дефолтной конфигурации
		raw := defaultHTTPServerConfig()

		// Применяем YAML конфигурацию поверх дефолтов
		if err := section.Unmarshal(raw); err == nil {
			return &httpServerConfig{Raw: *raw}, nil
		}
	}

	// Fallback на ENV
	return DefaultHTTPServerConfig()
}

func (s *httpServerConfig) ReadHeaderTimeout() time.Duration {
	return s.Raw.ReadHeaderTimeout
}

func (s *httpServerConfig) Address() string {
	return net.JoinHostPort(s.Raw.Host, strconv.Itoa(s.Raw.Port))
}

func (s *httpServerConfig) ReadTimeout() time.Duration { return s.Raw.ReadTimeout }

func (s *httpServerConfig) WriteTimeout() time.Duration { return s.Raw.WriteTimeout }

func (s *httpServerConfig) IdleTimeout() time.Duration { return s.Raw.IdleTimeout }

func (s *httpServerConfig) MaxHeaderBytes() int { return s.Raw.MaxHeaderBytes }

func (s *httpServerConfig) ShutdownTimeout() time.Duration { return s.Raw.ShutdownTimeout }

// HandlerTimeout returns application handler timeout; if not set, falls back to ReadTimeout
func (s *httpServerConfig) HandlerTimeout() time.Duration {
	if s.Raw.HandlerTimeout > 0 {
		return s.Raw.HandlerTimeout
	}
	return s.Raw.ReadTimeout
}
