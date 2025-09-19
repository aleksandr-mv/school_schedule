package permission

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/aleksandr-mv/school_schedule/rbac/internal/repository"
)

var _ def.PermissionRepository = (*permissionRepository)(nil)

type permissionRepository struct {
	readPool *pgxpool.Pool // Replica - только для чтения (SELECT)
}

func NewRepository(readPool *pgxpool.Pool) *permissionRepository {
	return &permissionRepository{readPool: readPool}
}
