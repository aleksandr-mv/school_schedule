package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/aleksandr-mv/school_schedule/iam/internal/repository"
)

var _ def.UserRepository = (*userRepository)(nil)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *userRepository {
	return &userRepository{pool: pool}
}
