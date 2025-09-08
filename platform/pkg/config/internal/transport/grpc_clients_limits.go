package transport

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// grpcClientLimitsRaw соответствует полям внутри секции grpc_client_limits.
type grpcClientLimitsRaw struct {
	MaxRecvMsgSize int           `mapstructure:"max_recv_msg_size" yaml:"max_recv_msg_size" env:"GRPC_MAX_REC_MSG_SIZE"`
	MaxSendMsgSize int           `mapstructure:"max_send_msg_size" yaml:"max_send_msg_size" env:"GRPC_MAX_SEND_MSG_SIZE"`
	Timeout        time.Duration `mapstructure:"timeout" yaml:"timeout" env:"GRPC_CLIENT_TIMEOUT"`
}

// grpcClientLimitsConfig хранит данные секции grpc_client_limits и реализует GRPCClientLimits.
type grpcClientLimitsConfig struct {
	Raw grpcClientLimitsRaw `yaml:"grpc_client_limits"`
}

var _ contracts.GRPCClientLimits = (*grpcClientLimitsConfig)(nil)

// defaultGRPCClientLimitsConfig возвращает конфигурацию gRPC client limits с дефолтными значениями
func defaultGRPCClientLimitsConfig() *grpcClientLimitsRaw {
	return &grpcClientLimitsRaw{
		MaxRecvMsgSize: 4 * 1024 * 1024,
		MaxSendMsgSize: 4 * 1024 * 1024,
		Timeout:        5 * time.Second,
	}
}

func DefaultGRPCClientLimitsConfig() (*grpcClientLimitsConfig, error) {
	raw := defaultGRPCClientLimitsConfig()

	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse gRPC client limits env: %w", err)
	}

	return &grpcClientLimitsConfig{Raw: *raw}, nil
}

// NewGRPCClientLimitsConfig создает конфигурацию лимитов gRPC клиента. YAML -> ENV -> валидация.
func NewGRPCClientLimitsConfig() (*grpcClientLimitsConfig, error) {
	if section := helpers.GetSection("grpc_client_limits"); section != nil {
		raw := defaultGRPCClientLimitsConfig()

		if err := section.Unmarshal(raw); err == nil {
			return &grpcClientLimitsConfig{Raw: *raw}, nil
		}
	}

	return DefaultGRPCClientLimitsConfig()
}

func (c *grpcClientLimitsConfig) MaxRecvMsgSize() int { return c.Raw.MaxRecvMsgSize }

func (c *grpcClientLimitsConfig) MaxSendMsgSize() int { return c.Raw.MaxSendMsgSize }

func (c *grpcClientLimitsConfig) Timeout() time.Duration { return c.Raw.Timeout }
