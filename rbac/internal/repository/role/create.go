package role

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

func (r *roleRepository) Create(ctx context.Context, name, description string) (uuid.UUID, error) {
	const query = `INSERT INTO roles (name, description) VALUES ($1, $2) RETURNING id`

	var id uuid.UUID
	if err := r.writePool.QueryRow(ctx, query, name, description).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return uuid.Nil, model.ErrRoleAlreadyExists
			}
		}
		return uuid.Nil, model.ErrFailedToCreateRole
	}

	return id, nil
}
