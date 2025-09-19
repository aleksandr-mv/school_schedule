package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/aleksandr-mv/school_schedule/iam/internal/repository"
)

var _ def.UserRepository = (*userRepository)(nil)

type userRepository struct {
	writePool *pgxpool.Pool
	readPool  *pgxpool.Pool
}

func NewRepository(writePool, readPool *pgxpool.Pool) *userRepository {
	return &userRepository{
		writePool: writePool,
		readPool:  readPool,
	}
}
