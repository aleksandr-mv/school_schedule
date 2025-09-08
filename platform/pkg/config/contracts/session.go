package contracts

import "time"

// SessionConfig определяет интерфейс для конфигурации сессий
type SessionConfig interface {
	// TTL возвращает время жизни сессии
	TTL() time.Duration
}
