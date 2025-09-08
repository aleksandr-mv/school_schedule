package model

import (
	"time"

	"github.com/google/uuid"
)

type UserCache struct {
	ID                  uuid.UUID
	Login               string
	Email               string
	CreatedAt           time.Time
	UpdatedAt           *time.Time
	NotificationMethods string
}
