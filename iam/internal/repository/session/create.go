package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	"github.com/aleksandr-mv/school_schedule/iam/internal/repository/converter"
)

func (r *sessionRepository) Create(ctx context.Context, user model.User, expiresAt time.Time) (uuid.UUID, error) {
	sessionID := uuid.New()
	now := time.Now()

	session := model.Session{
		ID:        sessionID,
		ExpiresAt: expiresAt,
		CreatedAt: now,
		UpdatedAt: now,
	}

	cacheKey := r.getCacheKey(sessionID.String())

	redisView, err := converter.CreateSessionCacheView(session, &user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: failed to create cache view: %w", model.ErrInvalidSessionData, err)
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return uuid.Nil, model.ErrSessionExpired
	}

	err = r.redis.HashSet(ctx, cacheKey, redisView)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %w", model.ErrFailedToStoreInCache, err)
	}

	if err = r.redis.Expire(ctx, cacheKey, ttl); err != nil {
		return uuid.Nil, fmt.Errorf("%w: failed to set expiration: %w", model.ErrFailedToStoreInCache, err)
	}

	return sessionID, nil
}
