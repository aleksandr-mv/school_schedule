package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository/converter"
	repoModel "github.com/aleksandr-mv/school_schedule/rbac/internal/repository/model"
)

func (r *roleRepository) Get(ctx context.Context, value string) (*model.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE id::text = $1 OR name = $1`

	rows, err := r.pool.Query(ctx, query, value)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	defer rows.Close()

	role, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[repoModel.Role])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to collect role: %w", err)
	}

	return converter.RoleToDomain(&role), nil
}
