package role_permission

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (r *rolePermissionRepository) Assign(ctx context.Context, roleID, permissionID string) error {
	query, args, err := sq.StatementBuilder.
		Insert("role_permissions").
		Columns("role_id", "permission_id").
		Values(roleID, permissionID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: failed to build insert query: %w", model.ErrInternal, err)
	}

	_, err = r.writePool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // Unique constraint violation
				return model.ErrPermissionAlreadyAssigned
			case "23503": // Foreign key constraint violation
				return fmt.Errorf("failed to assign permission: referenced entity not found")
			}
		}
		return fmt.Errorf("failed to assign permission to role: %w", err)
	}

	return nil
}
