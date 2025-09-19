package notification

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/aleksandr-mv/school_schedule/iam/internal/repository"
)

var _ def.NotificationRepository = (*notificationRepository)(nil)

type notificationRepository struct {
	writePool *pgxpool.Pool
	readPool  *pgxpool.Pool
}

func NewRepository(writePool, readPool *pgxpool.Pool) *notificationRepository {
	return &notificationRepository{
		writePool: writePool,
		readPool:  readPool,
	}
}
