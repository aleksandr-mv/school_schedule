package role_permission

import (
	"context"
	"fmt"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (r *rolePermissionRepository) Revoke(ctx context.Context, roleID, permissionID string) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`

	result, err := r.writePool.Exec(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to revoke permission from role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return model.ErrPermissionNotAssigned
	}

	return nil
}
