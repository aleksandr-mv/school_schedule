package user_role

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/aleksandr-mv/school_schedule/rbac/internal/repository"
)

var _ def.UserRoleRepository = (*userRoleRepository)(nil)

type userRoleRepository struct {
	writePool *pgxpool.Pool
	readPool  *pgxpool.Pool
}

func NewRepository(writePool, readPool *pgxpool.Pool) *userRoleRepository {
	return &userRoleRepository{
		writePool: writePool,
		readPool:  readPool,
	}
}
