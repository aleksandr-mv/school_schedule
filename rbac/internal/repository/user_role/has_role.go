package user_role

import (
	"context"
	"fmt"
)

func (r *userRoleRepository) HasRole(ctx context.Context, userID, roleID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_roles WHERE user_id = $1 AND role_id = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, userID, roleID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check role: %w", err)
	}

	return exists, nil
}
