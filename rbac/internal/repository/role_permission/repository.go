package role_permission

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/aleksandr-mv/school_schedule/rbac/internal/repository"
)

var _ def.RolePermissionRepository = (*rolePermissionRepository)(nil)

type rolePermissionRepository struct {
	writePool *pgxpool.Pool // Primary - для записи (INSERT, UPDATE, DELETE)
	readPool  *pgxpool.Pool // Replica - для чтения (SELECT)
}

func NewRepository(writePool, readPool *pgxpool.Pool) *rolePermissionRepository {
	return &rolePermissionRepository{
		writePool: writePool,
		readPool:  readPool,
	}
}
