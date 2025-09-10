package model

import (
	"time"

	"github.com/google/uuid"
)

// Role представляет роль пользователя
type Role struct {
	ID          uuid.UUID `validate:"required"`
	Name        string    `validate:"required,min=2,max=50"`
	Description string    `validate:"max=500"`
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func (r *Role) Validate() error {
	return validate.Struct(r)
}
