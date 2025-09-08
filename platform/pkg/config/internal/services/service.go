package services

import (
	"net"
	"strconv"
	"time"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// rawService соответствует одной записи services в YAML и env.
type rawService struct {
	Host    string        `mapstructure:"address" yaml:"address"`
	Port    int           `mapstructure:"port"    yaml:"port"`
	Timeout time.Duration `mapstructure:"timeout" yaml:"timeout"`
}

// serviceConfig хранит данные одной службы.
type serviceConfig struct {
	Raw rawService
}

func New() (contracts.ServicesConfig, error) {
	return newServicesConfig()
}

// defaultServiceConfig возвращает конфигурацию сервиса с дефолтными значениями
func defaultServiceConfig() rawService {
	return rawService{
		Host:    "",
		Port:    0,
		Timeout: 30 * time.Second, // только таймаут имеет смысл по умолчанию
	}
}

// newServiceConfig создаёт *serviceConfig из rawService.
func newServiceConfig(raw rawService) (*serviceConfig, error) {
	return &serviceConfig{Raw: raw}, nil
}

func (s *serviceConfig) Host() string { return s.Raw.Host }

func (s *serviceConfig) Port() int { return s.Raw.Port }

func (s *serviceConfig) Timeout() time.Duration { return s.Raw.Timeout }

func (s *serviceConfig) Address() string {
	return net.JoinHostPort(s.Raw.Host, strconv.Itoa(s.Raw.Port))
}
