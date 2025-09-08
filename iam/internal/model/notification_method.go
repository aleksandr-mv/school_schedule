package model

import (
	"time"

	"github.com/google/uuid"
)

type NotificationMethod struct {
	UserID       uuid.UUID  `json:"user_id" db:"user_id" validate:"required,uuid4"`
	ProviderName string     `json:"provider_name" db:"provider_name" validate:"required,max=100"`
	Target       string     `json:"target" db:"target" validate:"required,max=255"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func (nm *NotificationMethod) Validate() error {
	return validate.Struct(nm)
}
