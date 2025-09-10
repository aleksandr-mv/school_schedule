package role_permission

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/aleksandr-mv/school_schedule/rbac/internal/repository"
)

var _ def.RolePermissionRepository = (*rolePermissionRepository)(nil)

type rolePermissionRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *rolePermissionRepository {
	return &rolePermissionRepository{pool: pool}
}
