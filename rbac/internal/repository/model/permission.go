package model

import "github.com/google/uuid"

type Permission struct {
	ID       uuid.UUID `db:"id"`
	Resource string    `db:"resource"`
	Action   string    `db:"action"`
}
