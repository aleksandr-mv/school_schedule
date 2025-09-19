package user_role

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *userRoleRepository) GetRoleUsers(ctx context.Context, roleID string, limit int32, cursor string) ([]string, *string, error) {
	query := `
		SELECT user_id::text FROM user_roles
		WHERE role_id = $1 AND ($2::text = '' OR user_id > $2)
		ORDER BY user_id LIMIT $3`
	args := []interface{}{roleID, cursor, limit + 1}

	rows, err := r.readPool.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("get role users failed: %w", err)
	}
	defer rows.Close()

	userIDs, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, nil, fmt.Errorf("collect role users failed: %w", err)
	}

	var nextCursor *string
	if int32(len(userIDs)) > limit {
		next := userIDs[limit]
		userIDs = userIDs[:limit]
		nextCursor = &next
	}

	return userIDs, nextCursor, nil
}
