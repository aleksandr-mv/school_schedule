package user_role

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (r *userRoleRepository) AssignRole(ctx context.Context, userID, roleID string, assignedBy *string) error {
	queryBuilder := sq.StatementBuilder.
		Insert("user_roles").
		Columns("user_id", "role_id")

	if assignedBy != nil {
		queryBuilder = queryBuilder.Columns("assigned_by").Values(userID, roleID, *assignedBy)
	} else {
		queryBuilder = queryBuilder.Values(userID, roleID)
	}

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%w: failed to build insert query: %w", model.ErrInternal, err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // Unique constraint violation
				return model.ErrRoleAlreadyAssigned
			case "23503": // Foreign key constraint violation
				return fmt.Errorf("failed to assign role: referenced entity not found")
			}
		}
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}
