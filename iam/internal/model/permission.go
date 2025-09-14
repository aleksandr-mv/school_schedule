package model

import (
	"github.com/google/uuid"
)

// Permission представляет право доступа
type Permission struct {
	ID       uuid.UUID
	Resource string
	Action   string
}
