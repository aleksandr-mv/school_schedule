package user_role

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *userRoleRepository) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	query := `SELECT r.id::text 
		FROM roles r 
		JOIN user_roles ur ON r.id = ur.role_id 
		WHERE ur.user_id = $1
		ORDER BY r.name`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	roleIDs, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("failed to collect user role IDs: %w", err)
	}

	return roleIDs, nil
}
