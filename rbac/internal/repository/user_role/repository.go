package user_role

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/aleksandr-mv/school_schedule/rbac/internal/repository"
)

var _ def.UserRoleRepository = (*userRoleRepository)(nil)

type userRoleRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *userRoleRepository {
	return &userRoleRepository{pool: pool}
}
