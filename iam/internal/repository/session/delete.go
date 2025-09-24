package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
)

func (r *sessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) error {
	cacheKey := r.getCacheKey(sessionID.String())
	// DEL работает одинаково для обычных ключей и hash ключей
	if err := r.redis.Del(ctx, cacheKey); err != nil {
		return fmt.Errorf("%w: %w", model.ErrFailedToDeleteSession, err)
	}

	return nil
}
