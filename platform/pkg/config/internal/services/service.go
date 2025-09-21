package services

import (
	"net"
	"strconv"
	"time"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.ServiceConfig = (*ServiceConfig)(nil)

// rawServiceConfig для загрузки данных из YAML
type rawServiceConfig struct {
	Host    string        `mapstructure:"host"    yaml:"host"`
	Port    int           `mapstructure:"port"    yaml:"port"`
	Timeout time.Duration `mapstructure:"timeout" yaml:"timeout"`
}

// ServiceConfig публичная структура для использования
type ServiceConfig struct {
	raw rawServiceConfig
}

// defaultService возвращает rawServiceConfig с дефолтными значениями
func defaultService() rawServiceConfig {
	return rawServiceConfig{
		Host:    "",
		Port:    0,
		Timeout: 30 * time.Second,
	}
}

// newServiceConfig создает ServiceConfig из rawServiceConfig
func newServiceConfig(raw rawServiceConfig) *ServiceConfig {
	return &ServiceConfig{raw: raw}
}

// Методы для ServiceConfig интерфейса
func (s *ServiceConfig) Host() string           { return s.raw.Host }
func (s *ServiceConfig) Port() int              { return s.raw.Port }
func (s *ServiceConfig) Timeout() time.Duration { return s.raw.Timeout }

func (s *ServiceConfig) Address() string {
	return net.JoinHostPort(s.raw.Host, strconv.Itoa(s.raw.Port))
}
