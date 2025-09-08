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

// grpcRawConfig соответствует полям внутри секции grpc_server.
type grpcRawConfig struct {
	Host            string        `mapstructure:"host"               yaml:"host"                 env:"GRPC_HOST"`
	Port            int           `mapstructure:"port"               yaml:"port"                 env:"GRPC_PORT"`
	Timeout         time.Duration `mapstructure:"timeout"            yaml:"timeout"              env:"GRPC_TIMEOUT"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"       yaml:"idle_timeout"         env:"GRPC_IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"   yaml:"shutdown_timeout"     env:"GRPC_SHUTDOWN_TIMEOUT"`
}

// grpcServerConfig хранит данные секции grpc_server и реализует GRPCConfig.
type grpcServerConfig struct {
	Raw grpcRawConfig `yaml:"grpc_server"`
}

var _ contracts.GRPCServer = (*grpcServerConfig)(nil)

// defaultGRPCServerConfig возвращает конфигурацию gRPC server с дефолтными значениями
func defaultGRPCServerConfig() *grpcRawConfig {
	return &grpcRawConfig{
		Host:            "localhost",
		Port:            50055,
		Timeout:         30 * time.Second,
		IdleTimeout:     120 * time.Second,
		ShutdownTimeout: 10 * time.Second,
	}
}

// DefaultGRPCServerConfig читает конфиг gRPC из ENV
func DefaultGRPCServerConfig() (*grpcServerConfig, error) {
	raw := defaultGRPCServerConfig()
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse gRPC server env: %w", err)
	}

	return &grpcServerConfig{Raw: *raw}, nil
}

// NewGRPCServerConfig создает конфигурацию gRPC сервера. YAML -> ENV -> валидация.
func NewGRPCServerConfig() (*grpcServerConfig, error) {
	if section := helpers.GetSection("grpc_server"); section != nil {
		raw := defaultGRPCServerConfig()
		if err := section.Unmarshal(raw); err == nil {
			return &grpcServerConfig{Raw: *raw}, nil
		}
	}

	return DefaultGRPCServerConfig()
}

func (c *grpcServerConfig) Address() string {
	return net.JoinHostPort(c.Raw.Host, strconv.Itoa(c.Raw.Port))
}

func (c *grpcServerConfig) Timeout() time.Duration { return c.Raw.Timeout }

func (c *grpcServerConfig) IdleTimeout() time.Duration { return c.Raw.IdleTimeout }

func (c *grpcServerConfig) ShutdownTimeout() time.Duration { return c.Raw.ShutdownTimeout }
