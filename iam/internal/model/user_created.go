package model

import (
	"time"

	"github.com/google/uuid"
)

// UserCreated представляет событие создания пользователя для Kafka
type UserCreated struct {
	EventID   uuid.UUID `json:"event_id"`
	UserID    uuid.UUID `json:"user_id"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	RoleID    string    `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// NewUserCreated создает новое событие UserCreated
func NewUserCreated(user *User, roleID string) UserCreated {
	return UserCreated{
		EventID:   uuid.New(),
		UserID:    user.ID,
		Login:     user.Login,
		Email:     user.Email,
		RoleID:    roleID,
		CreatedAt: user.CreatedAt,
	}
}
