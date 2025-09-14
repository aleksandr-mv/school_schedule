package model

import (
	"time"

	"github.com/google/uuid"
)

// UserCreated представляет событие создания нового пользователя из IAM
type UserCreated struct {
	EventID   uuid.UUID `json:"event_id"`
	UserID    uuid.UUID `json:"user_id"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	RoleID    string    `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}
