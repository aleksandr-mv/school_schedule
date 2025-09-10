package model

import (
	"time"

	"github.com/google/uuid"
)

// RolePermission представляет связь роль-право
type RolePermission struct {
	RoleID       uuid.UUID `validate:"required"`
	PermissionID uuid.UUID `validate:"required"`
	CreatedAt    time.Time
}

func (rp *RolePermission) Validate() error {
	return validate.Struct(rp)
}
