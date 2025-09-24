package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/converter"
)

func (r *sessionRepository) Get(ctx context.Context, sessionID uuid.UUID) (*model.WhoAMI, error) {
	cacheKey := r.getCacheKey(sessionID.String())

	hash, err := r.redis.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, model.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: %w", model.ErrFailedToReadFromCache, err)
	}

	if len(hash) == 0 {
		return nil, model.ErrSessionNotFound
	}

	whoami, err := converter.FromRedisHash(hash)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to convert from hash: %w", model.ErrInvalidSessionData, err)
	}

	return whoami, nil
}
