package user_role

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository/converter"
	repoModel "github.com/aleksandr-mv/school_schedule/rbac/internal/repository/model"
)

func (r *userRoleRepository) GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error) {
	query := `SELECT r.id, r.name, r.description, r.created_at, r.updated_at 
		FROM roles r 
		JOIN user_roles ur ON r.id = ur.role_id 
		WHERE ur.user_id = $1
		ORDER BY r.name`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	roles, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[repoModel.Role])
	if err != nil {
		return nil, fmt.Errorf("failed to collect user roles: %w", err)
	}

	return converter.RolesToDomain(roles), nil
}
