package session

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/converter"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/model"
)

func (r *sessionRepository) Get(ctx context.Context, sessionID uuid.UUID) (*model.WhoAMI, error) {
	cacheKey := r.getCacheKey(sessionID.String())

	values, err := r.redis.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: %w", model.ErrFailedToReadFromCache, err)
	}

	if len(values) == 0 {
		return nil, model.ErrSessionNotFound
	}

	var sessionCacheView repoModel.SessionCacheView
	err = redigo.ScanStruct(values, &sessionCacheView)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to scan cache data: %w", model.ErrInvalidSessionData, err)
	}

	session, err := converter.SessionFromRedis(sessionCacheView.SessionRedisView)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to convert session data: %w", model.ErrInvalidSessionData, err)
	}

	user, err := converter.UserFromRedis(sessionCacheView.UserRedisView)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to convert user data: %w", model.ErrInvalidSessionData, err)
	}

	return &model.WhoAMI{
		Session: session,
		User:    *user,
	}, nil
}
