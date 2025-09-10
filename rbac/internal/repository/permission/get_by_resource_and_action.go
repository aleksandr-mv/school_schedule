package permission

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository/converter"
	repoModel "github.com/aleksandr-mv/school_schedule/rbac/internal/repository/model"
)

func (r *permissionRepository) GetByResourceAndAction(ctx context.Context, resource string, action string) (*model.Permission, error) {
	query := `SELECT id, resource, action FROM permissions WHERE resource = $1 AND action = $2`

	rows, err := r.pool.Query(ctx, query, resource, action)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission by resource and action: %w", err)
	}
	defer rows.Close()

	permission, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[repoModel.Permission])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("failed to collect permission by resource and action: %w", err)
	}

	return converter.PermissionToDomain(&permission), nil
}
