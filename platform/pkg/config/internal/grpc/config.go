package grpc

import (
	"net"
	"strconv"
	"time"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// rawConfig для загрузки данных из YAML/ENV
type rawConfig struct {
	// Настройки сервера
	Host            string        `mapstructure:"host"               yaml:"host"                 env:"GRPC_HOST"`
	Port            int           `mapstructure:"port"               yaml:"port"                 env:"GRPC_PORT"`
	Timeout         time.Duration `mapstructure:"timeout"            yaml:"timeout"              env:"GRPC_TIMEOUT"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"       yaml:"idle_timeout"         env:"GRPC_IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"   yaml:"shutdown_timeout"     env:"GRPC_SHUTDOWN_TIMEOUT"`

	// Настройки клиента
	MaxRecvMsgSize int           `mapstructure:"max_recv_msg_size" yaml:"max_recv_msg_size" env:"GRPC_MAX_REC_MSG_SIZE"`
	MaxSendMsgSize int           `mapstructure:"max_send_msg_size" yaml:"max_send_msg_size" env:"GRPC_MAX_SEND_MSG_SIZE"`
	ClientTimeout  time.Duration `mapstructure:"client_timeout"    yaml:"client_timeout"    env:"GRPC_CLIENT_TIMEOUT"`
}

// Config публичная структура для использования
type Config struct {
	raw rawConfig
}

// Компиляционная проверка
var _ contracts.GRPCConfig = (*Config)(nil)

// defaultConfig возвращает rawConfig с дефолтными значениями
func defaultConfig() rawConfig {
	return rawConfig{
		// Сервер
		Host:            "localhost",
		Port:            50051,
		Timeout:         30 * time.Second,
		IdleTimeout:     120 * time.Second,
		ShutdownTimeout: 10 * time.Second,

		// Клиент
		MaxRecvMsgSize: 4 * 1024 * 1024, // 4MB
		MaxSendMsgSize: 4 * 1024 * 1024, // 4MB
		ClientTimeout:  5 * time.Second,
	}
}

// Методы сервера
func (c *Config) Address() string                { return net.JoinHostPort(c.raw.Host, strconv.Itoa(c.raw.Port)) }
func (c *Config) Timeout() time.Duration         { return c.raw.Timeout }
func (c *Config) IdleTimeout() time.Duration     { return c.raw.IdleTimeout }
func (c *Config) ShutdownTimeout() time.Duration { return c.raw.ShutdownTimeout }

// Методы клиента
func (c *Config) MaxRecvMsgSize() int          { return c.raw.MaxRecvMsgSize }
func (c *Config) MaxSendMsgSize() int          { return c.raw.MaxSendMsgSize }
func (c *Config) ClientTimeout() time.Duration { return c.raw.ClientTimeout }
