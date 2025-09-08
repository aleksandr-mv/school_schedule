package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
)

func (r *sessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) error {
	cacheKey := r.getCacheKey(sessionID.String())
	if err := r.redis.Del(ctx, cacheKey); err != nil {
		return fmt.Errorf("%w: %w", model.ErrFailedToDeleteSession, err)
	}

	return nil
}
