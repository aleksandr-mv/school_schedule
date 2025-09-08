package postgres

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/viper"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

type rawConnectionConfig struct {
	Host            string `mapstructure:"host"               yaml:"host"               env:"DB_HOST"`
	Port            int    `mapstructure:"port"               yaml:"port"               env:"DB_PORT"`
	User            string `mapstructure:"user"               yaml:"user"               env:"POSTGRES_USER"`
	Password        string `mapstructure:"password"           yaml:"password"           env:"POSTGRES_PASSWORD"`
	Database        string `mapstructure:"database"           yaml:"database"           env:"POSTGRES_DB"`
	ApplicationName string `mapstructure:"application_name"   yaml:"application_name"   env:"DB_APPLICATION_NAME"`
}

type connectionConfig struct {
	Raw rawConnectionConfig `yaml:"connection"`
}

// Реализуем contracts.DBConnectionInfo
var _ contracts.DBConnection = (*connectionConfig)(nil)

func defaultConnectionConfig() *rawConnectionConfig {
	return &rawConnectionConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "password",
		Database:        "microservices",
		ApplicationName: "",
	}
}

// DefaultConnectionConfig читает конфиг соединения из ENV
func DefaultConnectionConfig() (*connectionConfig, error) {
	raw := defaultConnectionConfig()
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse connection config: %w", err)
	}
	return &connectionConfig{Raw: *raw}, nil
}

// NewConnectionConfig читает конфигурацию подключения (database.postgres.connection)
// из поддерева Viper. При отсутствии секции — fallback на переменные окружения.
func NewConnectionConfig(sub *viper.Viper) (*connectionConfig, error) {
	if sub == nil {
		return DefaultConnectionConfig()
	}
	connV := sub.Sub("connection")
	if connV == nil {
		return DefaultConnectionConfig()
	}
	raw := defaultConnectionConfig()
	if err := connV.Unmarshal(raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal connection config: %w", err)
	}
	return &connectionConfig{Raw: *raw}, nil
}

// URIWithTLS строит DSN с учётом TLS и application_name.
func (c *connectionConfig) URIWithTLS(tls contracts.TLSConfig) string {
	params := []string{}

	params = append(params, tls.BuildQueryParams()...)
	if c.Raw.ApplicationName != "" {
		params = append(params, "application_name="+c.Raw.ApplicationName)
	}

	query := ""
	if len(params) > 0 {
		query = "?" + strings.Join(params, "&")
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s%s",
		c.Raw.User, c.Raw.Password, c.Raw.Host, c.Raw.Port, c.Raw.Database, query,
	)
}

// Реализация интерфейса contracts.DBConnectionInfo
func (c *connectionConfig) Address() string {
	return net.JoinHostPort(c.Raw.Host, strconv.Itoa(c.Raw.Port))
}

func (c *connectionConfig) Database() string {
	return c.Raw.Database
}

func (c *connectionConfig) ApplicationName() string {
	return c.Raw.ApplicationName
}

func (c *connectionConfig) URI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.Raw.User, c.Raw.Password, c.Raw.Host, c.Raw.Port, c.Raw.Database,
	)
}
