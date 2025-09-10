package role_permission

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository/converter"
	repoModel "github.com/aleksandr-mv/school_schedule/rbac/internal/repository/model"
)

func (r *rolePermissionRepository) ListPermissionsByRole(ctx context.Context, value string) ([]*model.Permission, error) {
	query := `SELECT p.id, p.resource, p.action 
		FROM permissions p 
		JOIN role_permissions rp ON p.id = rp.permission_id 
		JOIN roles r ON rp.role_id = r.id
		WHERE r.id::text = $1 OR r.name = $1
		ORDER BY p.resource, p.action`

	rows, err := r.pool.Query(ctx, query, value)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	defer rows.Close()

	permissions, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[repoModel.Permission])
	if err != nil {
		return nil, fmt.Errorf("failed to collect role permissions: %w", err)
	}

	return converter.PermissionsToDomain(permissions), nil
}
