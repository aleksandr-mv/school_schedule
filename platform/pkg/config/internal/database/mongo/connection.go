package mongo

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/viper"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

type rawConnectionConfig struct {
	Host     string `mapstructure:"host"     yaml:"host"      env:"MONGO_HOST"`
	Port     int    `mapstructure:"port"     yaml:"port"      env:"MONGO_PORT"`
	User     string `mapstructure:"user"     yaml:"user"      env:"MONGO_INITDB_ROOT_USERNAME"`
	Password string `mapstructure:"password" yaml:"password"  env:"MONGO_INITDB_ROOT_PASSWORD"`
	Database string `mapstructure:"database" yaml:"database"  env:"MONGO_DATABASE"`
	AuthDB   string `mapstructure:"auth_db"  yaml:"auth_db"   env:"MONGO_AUTH_DB"`
}

// ConnectionConfig содержит параметры подключения к MongoDB.
type connectionConfig struct {
	Raw rawConnectionConfig
}

// Реализуем contracts.DBConnectionInfo
var _ contracts.DBConnection = (*connectionConfig)(nil)

// defaultConnectionConfig возвращает конфигурацию подключения MongoDB с дефолтными значениями
func defaultConnectionConfig() *rawConnectionConfig {
	return &rawConnectionConfig{
		Host:     "localhost",
		Port:     27017,
		User:     "root",
		Password: "password",
		Database: "microservices",
		AuthDB:   "admin",
	}
}

// DefaultConnectionConfig парсит env и возвращает ConnectionConfig.
func DefaultConnectionConfig() (*connectionConfig, error) {
	raw := defaultConnectionConfig()
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse DB config: %w", err)
	}
	return &connectionConfig{Raw: *raw}, nil
}

// NewConnectionConfig читает конфигурацию подключения MongoDB
// из поддерева Viper (ожидается database.mongo.connection, если sub = v.Sub("database.mongo")).
// При отсутствии секции — fallback на переменные окружения.
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
		return nil, fmt.Errorf("failed to unmarshal Mongo DB connection config: %w", err)
	}
	return &connectionConfig{Raw: *raw}, nil
}

// URIWithTLS строит MongoDB URI с учётом TLS параметров
func (c *connectionConfig) URIWithTLS(tls contracts.TLSConfig) string {
	params := []string{}
	if c.Raw.AuthDB != "" && c.Raw.AuthDB != c.Raw.Database {
		params = append(params, "authSource="+c.Raw.AuthDB)
	}

	params = append(params, tls.BuildQueryParams()...)

	query := ""
	if len(params) > 0 {
		query = "?" + strings.Join(params, "&")
	}

	cred := ""
	if c.Raw.User != "" && c.Raw.Password != "" {
		cred = fmt.Sprintf("%s:%s@", c.Raw.User, c.Raw.Password)
	}

	return fmt.Sprintf("mongodb://%s%s:%d/%s%s", cred, c.Raw.Host, c.Raw.Port, c.Raw.Database, query)
}

// URI для интерфейса DBConnectionInfo (без параметров TLS)
func (c *connectionConfig) URI() string {
	params := []string{}
	if c.Raw.AuthDB != "" && c.Raw.AuthDB != c.Raw.Database {
		params = append(params, "authSource="+c.Raw.AuthDB)
	}

	query := ""
	if len(params) > 0 {
		query = "?" + strings.Join(params, "&")
	}

	cred := ""
	if c.Raw.User != "" && c.Raw.Password != "" {
		cred = fmt.Sprintf("%s:%s@", c.Raw.User, c.Raw.Password)
	}

	return fmt.Sprintf("mongodb://%s%s:%d/%s%s", cred, c.Raw.Host, c.Raw.Port, c.Raw.Database, query)
}

// Реализация интерфейса contracts.DBConnectionInfo
func (c *connectionConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Raw.Host, c.Raw.Port)
}

func (c *connectionConfig) Database() string {
	return c.Raw.Database
}

func (c *connectionConfig) ApplicationName() string {
	return "" // MongoDB не использует application_name
}

// MongoDB-специфичные методы (для внутреннего использования)
func (c *connectionConfig) AuthDB() string { return c.Raw.AuthDB }

func (c *connectionConfig) String() string {
	return fmt.Sprintf(
		"ConnectionConfig{Host:%s, Port:%d, Database:%s, AuthDB:%s}",
		c.Raw.Host, c.Raw.Port, c.Raw.Database, c.Raw.AuthDB,
	)
}
