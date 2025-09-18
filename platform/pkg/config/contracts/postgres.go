package contracts

import "time"

// ============================================================================
// POSTGRESQL СПЕЦИФИЧНЫЕ ИНТЕРФЕЙСЫ
// ============================================================================

// PostgresConfig предоставляет доступ к PostgreSQL с поддержкой Read/Write Splitting.
//
// Поддерживает архитектуру Primary-Replica:
// - Primary: запись (INSERT, UPDATE, DELETE)
// - Replicas: чтение (SELECT) с балансировкой нагрузки
//
// Использование:
//
//	cfg := config.Postgres()
//	writePool := pgxpool.New(ctx, cfg.PrimaryDSN())
//	readPool := pgxpool.New(ctx, cfg.ReplicaDSN())
type PostgresConfig interface {
	// Primary соединение (запись)
	Primary() DBConnection
	PrimaryURI() string

	// Replica соединения (чтение)
	Replicas() []DBConnection
	ReplicaURI() string // возвращает случайную реплику

	// Общие настройки
	Pool() PostgresPoolConfig
}

// PostgresPoolConfig расширяет PoolConfig специфичными для PostgreSQL настройками.
type PostgresPoolConfig interface {
	PoolConfig
	MaxConnLifetime() time.Duration
	HealthCheckPeriod() time.Duration
}
