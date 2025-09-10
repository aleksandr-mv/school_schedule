package user_role

import (
	"context"
	"fmt"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (r *userRoleRepository) RevokeRole(ctx context.Context, userID, roleID string) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	result, err := r.pool.Exec(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to revoke role from user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return model.ErrRoleNotAssigned
	}

	return nil
}
