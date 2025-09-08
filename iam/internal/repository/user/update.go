package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	"github.com/aleksandr-mv/school_schedule/iam/internal/repository/converter"
	repoModel "github.com/aleksandr-mv/school_schedule/iam/internal/repository/model"
)

func (r *userRepository) Update(ctx context.Context, user model.User) (*model.User, error) {
	repoUser := converter.ToRepoUser(&user)
	builder := sq.StatementBuilder.
		Update("users").
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": repoUser.ID})

	if repoUser.Login != "" {
		builder = builder.Set("login", repoUser.Login)
	}
	if repoUser.Email != "" {
		builder = builder.Set("email", repoUser.Email)
	}
	if repoUser.PasswordHash != "" {
		builder = builder.Set("password_hash", repoUser.PasswordHash)
	}

	query, args, err := builder.
		Suffix("RETURNING id, login, email, password_hash, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", model.ErrInternal, err)
	}

	row := r.pool.QueryRow(ctx, query, args...)

	var u repoModel.User
	if err = row.Scan(&u.ID, &u.Login, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, r.mapDatabaseError(err, "update")
	}

	return converter.ToDomainUser(&u), nil
}
