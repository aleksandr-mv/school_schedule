package permission

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/converter"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/model"
)

func (r *permissionRepository) List(ctx context.Context) ([]*model.Permission, error) {
	query := `SELECT id, resource, action FROM permissions ORDER BY resource, action`

	rows, err := r.readPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	defer rows.Close()

	permissions, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[repoModel.Permission])
	if err != nil {
		return nil, fmt.Errorf("failed to collect permissions: %w", err)
	}

	return converter.PermissionsToDomain(permissions), nil
}
