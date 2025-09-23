package postgres

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.DBConnection = (*Connection)(nil)

// rawConnection для загрузки данных из YAML/ENV
type rawConnection struct {
	Host     string `mapstructure:"host"             yaml:"host"             env:"DB_HOST"`
	Port     int    `mapstructure:"port"             yaml:"port"             env:"DB_PORT"`
	User     string `mapstructure:"user"             yaml:"user"             env:"POSTGRES_USER"`
	Password string `mapstructure:"password"         yaml:"password"         env:"POSTGRES_PASSWORD"`
	DB       string `mapstructure:"database"         yaml:"database"         env:"POSTGRES_DB"`
	AppName  string `mapstructure:"application_name" yaml:"application_name" env:"DB_APPLICATION_NAME"`
}

// Connection публичная структура для использования
type Connection struct {
	raw rawConnection
}

// defaultConnection возвращает rawConnection с дефолтными значениями
func defaultConnection() rawConnection {
	return rawConnection{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		DB:       "postgres",
		AppName:  "app",
	}
}

// Методы для DBConnection интерфейса
func (conn *Connection) Address() string {
	return net.JoinHostPort(conn.raw.Host, strconv.Itoa(conn.raw.Port))
}
func (conn *Connection) Database() string        { return conn.raw.DB }
func (conn *Connection) ApplicationName() string { return conn.raw.AppName }

func (conn *Connection) DSN() string {
	// Правильная обработка аутентификации
	var auth string
	if conn.raw.Password != "" {
		auth = fmt.Sprintf("%s:%s", conn.raw.User, conn.raw.Password)
	} else {
		auth = conn.raw.User
	}

	// Параметры подключения
	params := []string{"sslmode=disable"}
	if conn.raw.AppName != "" {
		params = append(params, "application_name="+conn.raw.AppName)
	}

	query := ""
	if len(params) > 0 {
		query = "?" + strings.Join(params, "&")
	}

	// Правильный формат: postgresql:// (не postgres://)
	return fmt.Sprintf("postgresql://%s@%s:%d/%s%s",
		auth, conn.raw.Host, conn.raw.Port, conn.raw.DB, query)
}

func (conn *Connection) URI() string { return conn.DSN() } // alias
