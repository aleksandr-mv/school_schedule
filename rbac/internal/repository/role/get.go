package role

import (
	"context"
	"fmt"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/converter"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/model"
)

func (r *roleRepository) Get(ctx context.Context, id string) (*model.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at, deleted_at FROM roles WHERE id = $1 AND deleted_at IS NULL`

	row := r.readPool.QueryRow(ctx, query, id)

	var role repoModel.Role
	err := row.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt, &role.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return converter.RoleToDomain(&role), nil
}
