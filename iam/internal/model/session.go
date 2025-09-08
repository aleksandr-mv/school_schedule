package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `validate:"required,uuid4"`
	ExpiresAt time.Time `validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Session) Validate() error {
	return validate.Struct(s)
}
