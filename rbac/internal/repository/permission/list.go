package permission

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository/converter"
	repoModel "github.com/aleksandr-mv/school_schedule/rbac/internal/repository/model"
)

func (r *permissionRepository) List(ctx context.Context, filter *model.ListPermissionsFilter) ([]*model.Permission, error) {
	var query string
	var args []interface{}
	argIndex := 1

	baseQuery := `SELECT DISTINCT p.id, p.resource, p.action FROM permissions p`
	whereClauses := []string{}
	orderBy := ` ORDER BY p.resource, p.action`

	if filter.RoleID != nil {
		baseQuery = `SELECT DISTINCT p.id, p.resource, p.action FROM permissions p 
			JOIN role_permissions rp ON p.id = rp.permission_id`
		whereClauses = append(whereClauses, fmt.Sprintf("rp.role_id = $%d", argIndex))
		args = append(args, *filter.RoleID)
		argIndex++
	}

	if filter.Resource != nil && strings.TrimSpace(*filter.Resource) != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("LOWER(p.resource) ILIKE $%d", argIndex))
		args = append(args, "%"+strings.ToLower(strings.TrimSpace(*filter.Resource))+"%")
		argIndex++
	}

	if filter.Action != nil && strings.TrimSpace(*filter.Action) != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("LOWER(p.action) ILIKE $%d", argIndex))
		args = append(args, "%"+strings.ToLower(strings.TrimSpace(*filter.Action))+"%")
		argIndex++
	}

	if len(whereClauses) > 0 {
		query = baseQuery + " WHERE " + strings.Join(whereClauses, " AND ") + orderBy
	} else {
		query = baseQuery + orderBy
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	defer rows.Close()

	permissions, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[repoModel.Permission])
	if err != nil {
		return nil, fmt.Errorf("failed to collect permissions: %w", err)
	}

	return converter.PermissionsToDomain(permissions), nil
}
