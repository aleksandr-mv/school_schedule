package user_role

import (
	"context"
	"fmt"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (r *userRoleRepository) Revoke(ctx context.Context, userID, roleID string) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	result, err := r.writePool.Exec(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("%w: revoke role failed: %w", model.ErrInternal, err)
	}

	if n := result.RowsAffected(); n == 0 {
		return model.ErrRoleNotAssigned
	}

	return nil
}
