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
