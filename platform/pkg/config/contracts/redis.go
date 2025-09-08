package contracts

import "time"

// ============================================================================
// REDIS КОНФИГУРАЦИЯ
// ============================================================================

// RedisConfig описывает конфигурацию для работы с Redis.
// Включает настройки подключения, пула соединений и опциональные параметры.
// Конфигурация является опциональной - если Redis не настроен,
// IsEnabled() возвращает false.
type RedisConfig interface {
	// IsEnabled возвращает true, если Redis настроен
	IsEnabled() bool

	// Connection возвращает параметры подключения к Redis
	Connection() RedisConnection

	// Pool возвращает настройки пула соединений Redis
	Pool() RedisPoolConfig
}

// RedisConnection представляет параметры подключения к Redis.
type RedisConnection interface {
	// Host возвращает хост Redis сервера
	Host() string

	// Port возвращает порт Redis сервера
	Port() int

	// Password возвращает пароль для подключения (может быть пустым)
	Password() string

	// Database возвращает номер базы данных Redis (0-15)
	Database() int

	// Address возвращает полный адрес в формате "host:port"
	Address() string

	// URI возвращает строку подключения к Redis
	// Пример: "redis://:password@localhost:6379/0"
	URI() string
}

// RedisPoolConfig представляет настройки пула соединений Redis.
// Redis использует легкие соединения, поэтому настройки отличаются от PostgreSQL.
type RedisPoolConfig interface {
	// PoolSize возвращает максимальный размер пула соединений
	PoolSize() int

	// MinIdle возвращает минимальное количество простаивающих соединений
	MinIdle() int

	// MaxRetries возвращает максимальное количество попыток переподключения
	MaxRetries() int

	// DialTimeout возвращает таймаут установки соединения
	DialTimeout() time.Duration

	// ReadTimeout возвращает таймаут чтения из Redis
	ReadTimeout() time.Duration

	// WriteTimeout возвращает таймаут записи в Redis
	WriteTimeout() time.Duration

	// PoolTimeout возвращает таймаут ожидания соединения из пула
	PoolTimeout() time.Duration

	// IdleTimeout возвращает время, через которое неиспользуемое соединение считается устаревшим
	IdleTimeout() time.Duration

	// IdleCheckFreq возвращает частоту проверки простаивающих соединений
	IdleCheckFreq() time.Duration
}
