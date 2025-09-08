package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	UserID    uuid.UUID `db:"user_id"`
	SessionID uuid.UUID `db:"session_id"`
	CreatedAt time.Time `db:"created_at"`
}
