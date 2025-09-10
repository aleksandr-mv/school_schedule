package role

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository/converter"
	repoModel "github.com/aleksandr-mv/school_schedule/rbac/internal/repository/model"
)

func (r *roleRepository) List(ctx context.Context, nameFilter string) ([]*model.Role, error) {
	var query string
	var args []interface{}

	trimmedFilter := strings.TrimSpace(nameFilter)
	if len(trimmedFilter) > 0 {
		lowerFilter := strings.ToLower(trimmedFilter)
		query = `SELECT id, name, description, created_at, updated_at FROM roles WHERE LOWER(name) ILIKE $1 ORDER BY name ASC`
		args = []interface{}{"%" + lowerFilter + "%"}
	} else {
		query = `SELECT id, name, description, created_at, updated_at FROM roles ORDER BY name ASC`
		args = []interface{}{}
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	roles, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[repoModel.Role])
	if err != nil {
		return nil, fmt.Errorf("failed to collect roles: %w", err)
	}

	return converter.RolesToDomain(roles), nil
}
