package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/converter"
)

func (r *sessionRepository) Create(ctx context.Context, whoami *model.WhoAMI, expiresAt time.Time) (uuid.UUID, error) {
	sessionID := uuid.New()
	cacheKey := r.getCacheKey(sessionID.String())

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return uuid.Nil, model.ErrSessionExpired
	}

	hash, err := converter.ToRedisHash(whoami, sessionID, expiresAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: failed to convert to hash: %w", model.ErrInvalidSessionData, err)
	}

	err = r.redis.HSet(ctx, cacheKey, hash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: failed to store hash: %w", model.ErrFailedToStoreInCache, err)
	}

	err = r.redis.Expire(ctx, cacheKey, ttl)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: failed to set TTL: %w", model.ErrFailedToStoreInCache, err)
	}

	return sessionID, nil
}
