package role_permission

import (
	"context"
	"fmt"
)

func (r *rolePermissionRepository) HasPermission(ctx context.Context, roleID, permissionID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM role_permissions WHERE role_id = $1 AND permission_id = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, roleID, permissionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	return exists, nil
}
