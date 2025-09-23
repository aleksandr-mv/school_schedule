package role

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/converter"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/model"
)

func (r *roleRepository) List(ctx context.Context) ([]*model.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM roles 
		WHERE deleted_at IS NULL
		ORDER BY name ASC`

	rows, err := r.readPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	rawRows, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[repoModel.Role])
	if err != nil {
		return nil, fmt.Errorf("failed to collect roles: %w", err)
	}

	result, err := converter.RolesToDomain(rawRows)
	if err != nil {
		return nil, fmt.Errorf("failed to convert roles: %w", err)
	}

	return result, nil
}
