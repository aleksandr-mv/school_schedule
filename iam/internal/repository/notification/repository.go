package notification

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/aleksandr-mv/school_schedule/iam/internal/repository"
)

var _ def.NotificationRepository = (*notificationRepository)(nil)

type notificationRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *notificationRepository {
	return &notificationRepository{
		pool: pool,
	}
}
