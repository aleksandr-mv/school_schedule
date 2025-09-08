package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID
	Login               string `validate:"required,min=3,max=50"`
	Email               string `validate:"required,email,max=255"`
	PasswordHash        string `validate:"required"`
	NotificationMethods []*NotificationMethod
	CreatedAt           time.Time
	UpdatedAt           *time.Time
}

func (u *User) Validate() error {
	return validate.Struct(u)
}
