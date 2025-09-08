package model

import (
	"time"

	"github.com/google/uuid"
)

type NotificationMethod struct {
	UserID       uuid.UUID  `db:"user_id"`
	ProviderName string     `db:"provider_name"`
	Target       string     `db:"target"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}
