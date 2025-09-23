package user

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/converter"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/model"
)

func (r *userRepository) Get(ctx context.Context, value string) (*model.User, error) {
	query := `
		SELECT id, login, email, password_hash, created_at, updated_at
		FROM users
		WHERE id::text = $1 OR login = $1 OR email = $1`
	rows, err := r.readPool.Query(ctx, query, value)
	if err != nil {
		return nil, r.mapDatabaseError(err, "get")
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[repoModel.User])
	if err != nil {
		return nil, r.mapDatabaseError(err, "get")
	}

	return converter.ToDomainUser(&user), nil
}
