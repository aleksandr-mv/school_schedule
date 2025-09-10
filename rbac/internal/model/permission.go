package model

import (
	"github.com/google/uuid"
)

// Permission представляет право доступа
type Permission struct {
	ID       uuid.UUID
	Resource string `validate:"required,min=1,max=100"`
	Action   string `validate:"required,min=1,max=50"`
}

func (p *Permission) Validate() error {
	return validate.Struct(p)
}
