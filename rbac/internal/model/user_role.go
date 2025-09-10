package model

import (
	"time"

	"github.com/google/uuid"
)

// UserRole представляет связь пользователь-роль
type UserRole struct {
	UserID     uuid.UUID  `validate:"required"`
	RoleID     uuid.UUID  `validate:"required"`
	AssignedBy *uuid.UUID `validate:"required"`
	AssignedAt time.Time
}

func (ur *UserRole) Validate() error {
	return validate.Struct(ur)
}
