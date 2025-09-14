package model

import (
	"time"

	"github.com/google/uuid"
)

// UserRole представляет связь пользователь-роль
type UserRole struct {
	UserID     uuid.UUID
	RoleID     uuid.UUID
	AssignedBy *uuid.UUID
	AssignedAt time.Time
}
