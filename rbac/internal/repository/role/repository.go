package role

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/aleksandr-mv/school_schedule/rbac/internal/repository"
)

var _ def.RoleRepository = (*roleRepository)(nil)

type roleRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *roleRepository {
	return &roleRepository{pool: pool}
}
