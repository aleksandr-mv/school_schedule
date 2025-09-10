package role

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository/converter"
)

func (r *roleRepository) Update(ctx context.Context, updateRole *model.UpdateRole) error {
	updates, err := converter.UpdateRoleToRepo(updateRole.Name, updateRole.Description)
	if err != nil {
		return fmt.Errorf("failed to prepare update data: %w", err)
	}

	queryBuilder := sq.StatementBuilder.
		Update("roles").
		Set("updated_at", "NOW()").
		Where(sq.Eq{"id": updateRole.ID})

	for key, value := range updates {
		queryBuilder = queryBuilder.Set(key, value)
	}

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%w: failed to build update query: %w", model.ErrInternal, err)
	}

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // Unique constraint violation
				return model.ErrRoleAlreadyExists
			case "23503": // Foreign key constraint violation
				return fmt.Errorf("failed to update role: referenced entity not found")
			}
		}
		return fmt.Errorf("failed to update role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return model.ErrRoleNotFound
	}

	return nil
}
