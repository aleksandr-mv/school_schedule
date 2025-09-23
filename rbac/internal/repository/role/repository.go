package role

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository"
)

var _ def.RoleRepository = (*roleRepository)(nil)

type roleRepository struct {
	writePool *pgxpool.Pool // Primary - для записи (INSERT, UPDATE, DELETE)
	readPool  *pgxpool.Pool // Replica - для чтения (SELECT)
}

func NewRepository(writePool, readPool *pgxpool.Pool) *roleRepository {
	return &roleRepository{
		writePool: writePool,
		readPool:  readPool,
	}
}
