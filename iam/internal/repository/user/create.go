package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/converter"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/model"
)

func (r *userRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	user.ID = uuid.New()
	repoUser := converter.ToRepoUser(&user)

	query, args, err := sq.StatementBuilder.
		Insert("users").
		Columns("id", "login", "email", "password_hash").
		Values(repoUser.ID, repoUser.Login, repoUser.Email, repoUser.PasswordHash).
		Suffix("RETURNING id, login, email, password_hash, created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to build insert query: %w", model.ErrInternal, err)
	}

	rows, err := r.writePool.Query(ctx, query, args...)
	if err != nil {
		return nil, r.mapDatabaseError(err, "create")
	}
	defer rows.Close()

	createdUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[repoModel.User])
	if err != nil {
		return nil, r.mapDatabaseError(err, "create")
	}

	return converter.ToDomainUser(&createdUser), nil
}
