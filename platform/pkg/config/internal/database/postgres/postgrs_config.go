// Package postgres предоставляет конфигурацию для подключения к PostgreSQL.
//
// Этот пакет инкапсулирует все аспекты конфигурации PostgreSQL:
// параметры подключения, настройки пула соединений и TLS конфигурацию.
// Поддерживает загрузку из YAML файлов с fallback на ENV переменные.
//
// Основные компоненты:
//   - Connection: хост, порт, база данных, учетные данные
//   - Pool: размер пула, таймауты, ограничения соединений
//   - TLS: настройки SSL/TLS шифрования
//
// Пример YAML конфигурации:
//
//	database:
//	  postgres:
//	    connection:
//	      host: localhost
//	      port: 5432
//	      database: mydb
//	      user: postgres
//	      password: secret
//	    pool:
//	      max_connections: 10
//	      min_connections: 2
//	      connect_timeout: 30s
//	    tls:
//	      enabled: true
//	      ssl_mode: require
package postgres

import (
	"fmt"

	contracts "github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// ============================================================================
// СТРУКТУРЫ КОНФИГУРАЦИИ POSTGRESQL
// ============================================================================

// rawPostgresConfig представляет сырую структуру конфигурации PostgreSQL из YAML.
// Объединяет все компоненты: подключение, пул соединений и TLS настройки.
type rawPostgresConfig struct {
	TLS        *tlsConfig        `yaml:"tls"`        // настройки TLS/SSL
	Connection *connectionConfig `yaml:"connection"` // параметры подключения
	Pool       *poolConfig       `yaml:"pool"`       // настройки пула соединений
}

// postgresConfig реализует contracts.PostgresConfig и предоставляет
// валидированную конфигурацию PostgreSQL для использования в приложении.
type postgresConfig struct {
	Raw rawPostgresConfig `yaml:"postgres"`
}

// Компиляционная проверка реализации интерфейса
var _ contracts.PostgresConfig = (*postgresConfig)(nil)

// ============================================================================
// КОНСТРУКТОРЫ И ФАБРИЧНЫЕ ФУНКЦИИ
// ============================================================================

// NewPostgresConfig создает и валидирует полную конфигурацию PostgreSQL.
//
// Функция пытается загрузить конфигурацию из YAML секции "database.postgres".
// При отсутствии или ошибке чтения YAML использует fallback на ENV переменные.
// Все компоненты (TLS, Connection, Pool) создаются и валидируются независимо.
//
// Порядок загрузки:
//  1. Чтение секции "database.postgres" из YAML
//  2. Создание и валидация TLS конфигурации
//  3. Создание и валидация Connection конфигурации
//  4. Создание и валидация Pool конфигурации
//  5. Fallback на ENV переменные при ошибках
//
// Возвращает:
//   - *postgresConfig: валидированную конфигурацию PostgreSQL
//   - error: при критических ошибках валидации
func NewPostgresConfig() (*postgresConfig, error) {
	if psql := helpers.GetSection("database.postgres"); psql != nil {
		tlsCfg, errTLS := NewTLSConfig(psql)
		if errTLS != nil {
			return nil, fmt.Errorf("tls config: %w", errTLS)
		}
		dbCfg, errDB := NewConnectionConfig(psql)
		if errDB != nil {
			return nil, fmt.Errorf("connection config: %w", errDB)
		}
		poolCfg, errPool := NewPoolConfig(psql)
		if errPool != nil {
			return nil, fmt.Errorf("pool config: %w", errPool)
		}

		return &postgresConfig{Raw: rawPostgresConfig{TLS: tlsCfg, Connection: dbCfg, Pool: poolCfg}}, nil
	}

	tlsCfg, err := DefaultTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("tls config: %w", err)
	}

	dbCfg, err := DefaultConnectionConfig()
	if err != nil {
		return nil, fmt.Errorf("connection config: %w", err)
	}

	poolCfg, err := DefaultPoolConfig()
	if err != nil {
		return nil, fmt.Errorf("pool config: %w", err)
	}

	return &postgresConfig{Raw: rawPostgresConfig{TLS: tlsCfg, Connection: dbCfg, Pool: poolCfg}}, nil
}

// ============================================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА contracts.PostgresConfig
// ============================================================================

// URI возвращает полную строку подключения PostgreSQL с TLS параметрами.
//
// Метод объединяет базовые параметры подключения с TLS настройками для
// создания готовой URI строки, которую можно использовать с pgxpool.
//
// Возвращает строку формата:
//
//	"postgresql://user:password@host:port/database?sslmode=require&..."
func (c *postgresConfig) URI() string {
	return c.Raw.Connection.URIWithTLS(c.Raw.TLS)
}

// TLS возвращает конфигурацию TLS/SSL шифрования.
//
// Предоставляет доступ к настройкам безопасности соединения:
// режим SSL, сертификаты, проверка хоста и другие TLS параметры.
func (c *postgresConfig) TLS() contracts.TLSConfig { return c.Raw.TLS }

// Database возвращает базовые параметры подключения к базе данных.
//
// Предоставляет унифицированный доступ к основным параметрам:
// адрес сервера, имя базы данных, учетные данные и имя приложения.
func (c *postgresConfig) Database() contracts.DBConnection { return c.Raw.Connection }

// Pool возвращает настройки пула соединений PostgreSQL.
//
// Предоставляет доступ к специфичным для PostgreSQL настройкам пула:
// максимальное и минимальное количество соединений, таймауты подключения и закрытия.
func (c *postgresConfig) Pool() contracts.PostgresPoolConfig { return c.Raw.Pool }
