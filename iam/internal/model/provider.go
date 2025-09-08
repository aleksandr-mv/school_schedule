package model

import (
	"time"
)

type Provider struct {
	ID          int        `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" validate:"required,max=100"`
	Description *string    `json:"description,omitempty" db:"description"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

func (p *Provider) Validate() error {
	return validate.Struct(p)
}
