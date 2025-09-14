package user_role

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (r *userRoleRepository) Assign(ctx context.Context, userID, roleID, assignedBy string) error {
	query := `INSERT INTO user_roles (user_id, role_id, assigned_by) VALUES ($1, $2, $3)`
	result, err := r.pool.Exec(ctx, query, userID, roleID, assignedBy)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return model.ErrRoleAlreadyAssigned
			case "23503": // foreign_key_violation
				return model.ErrUserRoleNotFound
			}
		}
		return fmt.Errorf("%w: assign role failed: %w", model.ErrInternal, err)
	}

	if n := result.RowsAffected(); n == 0 {
		return fmt.Errorf("%w: assign role failed: no rows affected", model.ErrInternal)
	}

	return nil
}
