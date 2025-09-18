package contracts

import "time"

// ============================================================================
// MONGODB СПЕЦИФИЧНЫЕ ИНТЕРФЕЙСЫ
// ============================================================================

// MongoConfig предоставляет доступ к полной конфигурации MongoDB с Primary-Replica архитектурой.
//
// Поддерживает разделение нагрузки между primary (запись) и replicas (чтение).
// MongoDB драйвер поддерживает специфичные настройки пула, отличные от PostgreSQL.
//
// Primary-Replica методы:
//   - Primary(): primary соединение для записи
//   - PrimaryURI(): URI для записи
//   - Replicas(): все replica соединения для чтения
//   - ReplicaURI(): случайная replica URI для чтения
//   - Pool(): настройки пула соединений
//
// Обратная совместимость:
//   - Database(): возвращает primary connection
//   - URI(): возвращает primary URI
//
// Использование:
//
//	cfg := config.Mongo()
//	if cfg != nil {
//		// Для записи
//		writeClient := mongo.Connect(ctx, options.Client().ApplyURI(cfg.PrimaryURI()))
//		// Для чтения
//		readClient := mongo.Connect(ctx, options.Client().ApplyURI(cfg.ReplicaURI()))
//	}
type MongoConfig interface {
	Primary() DBConnection
	PrimaryURI() string
	Replicas() []DBConnection
	ReplicaURI() string
	Pool() MongoPoolConfig
}

// MongoPoolConfig расширяет PoolConfig специфичными для MongoDB настройками.
type MongoPoolConfig interface {
	PoolConfig
	MaxConnecting() int
	HeartbeatPeriod() time.Duration
	ServerSelectionTimeout() time.Duration
}
