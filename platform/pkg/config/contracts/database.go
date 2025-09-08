package contracts

import "time"

// ============================================================================
// ОБЩИЕ ИНТЕРФЕЙСЫ ДЛЯ БАЗ ДАННЫХ
// ============================================================================

// DBConnection представляет унифицированный интерфейс для подключения к любой БД.
//
// Этот интерфейс обеспечивает единообразный доступ к параметрам подключения
// для разных типов баз данных (PostgreSQL, MongoDB). Каждая реализация
// адаптирует специфичные для драйвера настройки под общий API.
//
// Методы:
//   - Address(): возвращает адрес в формате "host:port"
//   - Database(): возвращает имя базы данных
//   - ApplicationName(): возвращает имя приложения (может быть пустым для MongoDB)
//   - URI(): возвращает готовую строку подключения без TLS параметров
//
// Примеры возвращаемых значений:
//
//	PostgreSQL: "postgresql://user:pass@localhost:5432/dbname"
//	MongoDB: "mongodb://user:pass@localhost:27017/dbname"
type DBConnection interface {
	Address() string         // host:port
	Database() string        // имя базы данных
	ApplicationName() string // имя приложения (для MongoDB может быть пустым)
	URI() string             // готовая строка подключения
}

// PoolConfig представляет общие настройки пула соединений.
type PoolConfig interface {
	MaxSize() int
	MinSize() int
	MaxIdleTime() time.Duration
	ConnectTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

// DatabaseConfig объединяет методы для работы с БД.
type DatabaseConfig interface {
	HasPostgres() bool
	HasMongo() bool
	PostgresDSN() string
	MongoURI() string
	PostgresConnection() (DBConnection, bool)
	MongoConnection() (DBConnection, bool)
	PostgresPool() (PostgresPoolConfig, bool)
	MongoPool() (MongoPoolConfig, bool)
}

// ============================================================================
// POSTGRESQL СПЕЦИФИЧНЫЕ ИНТЕРФЕЙСЫ
// ============================================================================

// PostgresConfig предоставляет доступ к полной конфигурации PostgreSQL.
//
// Этот интерфейс объединяет все аспекты конфигурации PostgreSQL:
// параметры подключения, настройки пула соединений и TLS конфигурацию.
//
// Методы:
//   - TLS(): настройки TLS/SSL для шифрования соединения
//   - Database(): базовые параметры подключения (хост, порт, БД)
//   - Pool(): настройки пула соединений (размер, таймауты)
//   - URI(): полная строка подключения с TLS параметрами
//
// Использование:
//
//	cfg := config.DatabaseConfig().PostgresConnection()
//	pool := pgxpool.New(ctx, cfg.URI())
type PostgresConfig interface {
	TLS() TLSConfig
	Database() DBConnection
	Pool() PostgresPoolConfig
	URI() string
}

// PostgresPoolConfig расширяет PoolConfig специфичными для PostgreSQL настройками.
type PostgresPoolConfig interface {
	PoolConfig
	MaxConnLifetime() time.Duration
	HealthCheckPeriod() time.Duration
}

// ============================================================================
// MONGODB СПЕЦИФИЧНЫЕ ИНТЕРФЕЙСЫ
// ============================================================================

// MongoConfig предоставляет доступ к полной конфигурации MongoDB.
//
// Этот интерфейс объединяет все аспекты конфигурации MongoDB:
// параметры подключения, настройки пула соединений и TLS конфигурацию.
// MongoDB драйвер поддерживает специфичные настройки пула, отличные от PostgreSQL.
//
// Методы:
//   - TLS(): настройки TLS/SSL для шифрования соединения
//   - Database(): базовые параметры подключения (хост, порт, БД)
//   - Pool(): настройки пула соединений (размер, таймауты, MongoDB-специфичные)
//   - URI(): полная строка подключения с TLS параметрами
//
// Использование:
//
//	cfg := config.DatabaseConfig().MongoConnection()
//	client := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI()))
type MongoConfig interface {
	TLS() TLSConfig
	Database() DBConnection
	Pool() MongoPoolConfig
	URI() string
}

// MongoPoolConfig расширяет PoolConfig специфичными для MongoDB настройками.
type MongoPoolConfig interface {
	PoolConfig
	MaxConnecting() int
}
