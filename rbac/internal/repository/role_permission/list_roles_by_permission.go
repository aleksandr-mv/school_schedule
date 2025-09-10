package role_permission

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository/converter"
	repoModel "github.com/aleksandr-mv/school_schedule/rbac/internal/repository/model"
)

func (r *rolePermissionRepository) ListRolesByPermission(ctx context.Context, value string) ([]*model.Role, error) {
	query := `SELECT r.id, r.name, r.description, r.created_at, r.updated_at 
		FROM roles r 
		JOIN role_permissions rp ON r.id = rp.role_id 
		JOIN permissions p ON rp.permission_id = p.id
		WHERE p.id::text = $1 OR (p.resource = $1 AND p.action = $1)
		ORDER BY r.name`

	rows, err := r.pool.Query(ctx, query, value)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles by permission: %w", err)
	}
	defer rows.Close()

	roles, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[repoModel.Role])
	if err != nil {
		return nil, fmt.Errorf("failed to collect roles by permission: %w", err)
	}

	return converter.RolesToDomain(roles), nil
}
