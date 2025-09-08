// Package mongo предоставляет конфигурацию для подключения к MongoDB.
//
// Этот пакет инкапсулирует все аспекты конфигурации MongoDB:
// параметры подключения, настройки пула соединений и TLS конфигурацию.
// Поддерживает загрузку из YAML файлов с fallback на ENV переменные.
//
// MongoDB отличается от PostgreSQL специфичными настройками пула:
//   - MaxPoolSize: максимальное количество соединений в пуле
//   - MinPoolSize: минимальное количество активных соединений
//   - MaxConnecting: максимальное количество одновременных подключений
//   - ConnectTimeout/SocketTimeout/MaxIdleTime: различные таймауты
//
// Пример YAML конфигурации:
//
//	database:
//	  mongo:
//	    connection:
//	      host: localhost
//	      port: 27017
//	      database: mydb
//	      user: mongoadmin
//	      password: secret
//	      auth_db: admin
//	    pool:
//	      max_pool_size: 10
//	      min_pool_size: 2
//	      max_connecting: 5
//	      connect_timeout: 30s
//	    tls:
//	      enabled: true
//	      ssl_mode: require
package mongo

import (
	"fmt"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// ============================================================================
// СТРУКТУРЫ КОНФИГУРАЦИИ MONGODB
// ============================================================================

// rawMongoConfig представляет сырую структуру конфигурации MongoDB из YAML.
// Объединяет все компоненты: подключение, пул соединений и TLS настройки.
type rawMongoConfig struct {
	TLS        *tlsConfig        // настройки TLS/SSL
	Connection *connectionConfig // параметры подключения к MongoDB
	Pool       *poolConfig       // настройки пула соединений (MongoDB-специфичные)
}

// mongoConfig реализует contracts.MongoConfig и предоставляет
// валидированную конфигурацию MongoDB для использования в приложении.
type mongoConfig struct {
	raw rawMongoConfig
}

// Компиляционная проверка реализации интерфейса
var _ contracts.MongoConfig = (*mongoConfig)(nil)

// ============================================================================
// КОНСТРУКТОРЫ И ФАБРИЧНЫЕ ФУНКЦИИ
// ============================================================================

// NewMongoConfig создает и валидирует полную конфигурацию MongoDB.
//
// Функция пытается загрузить конфигурацию из YAML секции "database.mongo".
// При отсутствии или ошибке чтения YAML использует fallback на ENV переменные.
// Все компоненты (TLS, Connection, Pool) создаются и валидируются независимо.
//
// Особенности MongoDB конфигурации:
//   - Поддержка auth_db для аутентификации
//   - MongoDB-специфичные настройки пула соединений
//   - Адаптация параметров подключения под mongo-driver
//
// Порядок загрузки:
//  1. Чтение секции "database.mongo" из YAML
//  2. Создание и валидация TLS конфигурации
//  3. Создание и валидация Connection конфигурации
//  4. Создание и валидация Pool конфигурации
//  5. Fallback на ENV переменные при ошибках
//
// Возвращает:
//   - *mongoConfig: валидированную конфигурацию MongoDB
//   - error: при критических ошибках валидации
func NewMongoConfig() (*mongoConfig, error) {
	if mg := helpers.GetSection("database.mongo"); mg != nil {
		tlsCfg, err := NewTLSConfig(mg)
		if err != nil {
			return nil, fmt.Errorf("tls config: %w", err)
		}

		dbCfg, err := NewConnectionConfig(mg)
		if err != nil {
			return nil, fmt.Errorf("connection config: %w", err)
		}

		poolCfg, err := NewPoolConfig(mg)
		if err != nil {
			return nil, fmt.Errorf("pool config: %w", err)
		}

		return &mongoConfig{raw: rawMongoConfig{TLS: tlsCfg, Connection: dbCfg, Pool: poolCfg}}, nil
	}

	// Fallback на ENV
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

	return &mongoConfig{raw: rawMongoConfig{TLS: tlsCfg, Connection: dbCfg, Pool: poolCfg}}, nil
}

// ============================================================================
// РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА contracts.MongoConfig
// ============================================================================

// URI возвращает полную строку подключения MongoDB с TLS параметрами.
//
// Метод объединяет базовые параметры подключения с TLS настройками для
// создания готовой URI строки, которую можно использовать с mongo-driver.
//
// Возвращает строку формата:
//
//	"mongodb://user:password@host:port/database?authSource=admin&ssl=true&..."
func (c *mongoConfig) URI() string {
	return c.raw.Connection.URIWithTLS(c.raw.TLS)
}

// TLS возвращает конфигурацию TLS/SSL шифрования.
//
// Предоставляет доступ к настройкам безопасности соединения:
// режим SSL, сертификаты, проверка хоста и другие TLS параметры.
func (c *mongoConfig) TLS() contracts.TLSConfig {
	return c.raw.TLS
}

// Database возвращает базовые параметры подключения к базе данных.
//
// Предоставляет унифицированный доступ к основным параметрам:
// адрес сервера, имя базы данных, учетные данные. Для MongoDB
// ApplicationName() может возвращать пустую строку.
func (c *mongoConfig) Database() contracts.DBConnection {
	return c.raw.Connection
}

// Pool возвращает настройки пула соединений MongoDB.
//
// Предоставляет доступ к специфичным для MongoDB настройкам пула:
// MaxPoolSize, MinPoolSize, MaxConnecting и различные таймауты,
// которые отличаются от PostgreSQL настроек.
func (c *mongoConfig) Pool() contracts.MongoPoolConfig {
	return c.raw.Pool
}
