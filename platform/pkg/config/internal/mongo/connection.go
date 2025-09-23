package mongo

import (
	"fmt"
	"net"
	"strconv"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.DBConnection = (*Connection)(nil)

// rawConnection для загрузки данных из YAML/ENV
type rawConnection struct {
	Host     string `mapstructure:"host"     yaml:"host"     env:"MONGO_HOST"`
	Port     int    `mapstructure:"port"     yaml:"port"     env:"MONGO_PORT"`
	User     string `mapstructure:"user"     yaml:"user"     env:"MONGO_USER"`
	Password string `mapstructure:"password" yaml:"password" env:"MONGO_PASSWORD"`
	DB       string `mapstructure:"database" yaml:"database" env:"MONGO_DB"`
	AuthDB   string `mapstructure:"auth_db"  yaml:"auth_db"  env:"MONGO_AUTH_DB"`
}

// Connection публичная структура для использования
type Connection struct {
	raw rawConnection
}

// defaultConnection возвращает rawConnection с дефолтными значениями
func defaultConnection() rawConnection {
	return rawConnection{
		Host:     "localhost",
		Port:     27017,
		User:     "",
		Password: "",
		DB:       "admin",
		AuthDB:   "admin",
	}
}

// Методы для DBConnection интерфейса
func (conn *Connection) Address() string {
	return net.JoinHostPort(conn.raw.Host, strconv.Itoa(conn.raw.Port))
}
func (conn *Connection) Database() string        { return conn.raw.DB }
func (conn *Connection) ApplicationName() string { return "" } // MongoDB не использует application_name

func (conn *Connection) URI() string {
	cred := ""
	if conn.raw.User != "" && conn.raw.Password != "" {
		cred = fmt.Sprintf("%s:%s@", conn.raw.User, conn.raw.Password)
	}

	authSource := ""
	if conn.raw.AuthDB != "" && conn.raw.AuthDB != conn.raw.DB {
		authSource = "?authSource=" + conn.raw.AuthDB
	}

	return fmt.Sprintf("mongodb://%s%s:%d/%s%s",
		cred, conn.raw.Host, conn.raw.Port, conn.raw.DB, authSource)
}
