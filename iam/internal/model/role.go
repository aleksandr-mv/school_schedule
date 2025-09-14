package model

import (
	"time"

	"github.com/google/uuid"
)

// Role представляет роль пользователя
type Role struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}
