package session

import (
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/cache"
)

var _ repository.SessionRepository = (*sessionRepository)(nil)

type sessionRepository struct {
	redis cache.RedisClient
}

func NewRepository(redis cache.RedisClient) *sessionRepository {
	return &sessionRepository{
		redis: redis,
	}
}
