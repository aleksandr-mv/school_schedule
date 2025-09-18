package contracts

import "time"

// ============================================================================
// REDIS КЛАСТЕР КОНФИГУРАЦИЯ
// ============================================================================

// RedisConfig описывает конфигурацию для работы с Redis Cluster.
// Только кластерный режим, без обратной совместимости.
type RedisConfig interface {
	// Cluster возвращает конфигурацию кластера
	Cluster() RedisClusterConfig

	// Pool возвращает настройки пула соединений Redis
	Pool() RedisPoolConfig
}

// RedisClusterConfig представляет конфигурацию Redis кластера.
type RedisClusterConfig interface {
	// IsEnabled возвращает true, если кластер настроен (есть узлы)
	IsEnabled() bool

	// Nodes возвращает список узлов кластера
	Nodes() []string

	// Password возвращает пароль для кластера (может быть пустым)
	Password() string

	// MaxRedirects максимальное количество редиректов
	MaxRedirects() int

	// ReadOnlyCommands позволяет читать с реплик
	ReadOnlyCommands() bool

	// RouteByLatency маршрутизация по задержке
	RouteByLatency() bool

	// RouteRandomly случайная маршрутизация
	RouteRandomly() bool

	// NodesAddresses возвращает узлы через запятую для подключения
	NodesAddresses() string
}

// RedisPoolConfig представляет настройки пула соединений Redis.
type RedisPoolConfig interface {
	// MaxActive максимальное количество активных соединений
	MaxActive() int

	// MaxIdle максимальное количество простаивающих соединений
	MaxIdle() int

	// IdleTimeout время жизни простаивающего соединения
	IdleTimeout() time.Duration

	// ConnTimeout таймаут установки соединения
	ConnTimeout() time.Duration

	// ReadTimeout таймаут чтения
	ReadTimeout() time.Duration

	// WriteTimeout таймаут записи
	WriteTimeout() time.Duration

	// PoolTimeout таймаут получения соединения из пула
	PoolTimeout() time.Duration
}
