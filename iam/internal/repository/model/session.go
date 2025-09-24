package model

// SessionRedisView представляет сессию в Redis
type SessionRedisView struct {
	ID          string `redis:"session_id"`
	ExpiresAtNs int64  `redis:"session_expires_at"`
	CreatedAtNs int64  `redis:"session_created_at"`
	UpdatedAtNs int64  `redis:"session_updated_at"`
}
