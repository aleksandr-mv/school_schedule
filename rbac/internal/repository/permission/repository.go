package permission

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/aleksandr-mv/school_schedule/rbac/internal/repository"
)

var _ def.PermissionRepository = (*permissionRepository)(nil)

type permissionRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *permissionRepository {
	return &permissionRepository{pool: pool}
}
