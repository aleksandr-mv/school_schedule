package enriched_role

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/converter"
)

func (r *repository) Set(ctx context.Context, role *model.EnrichedRole, expiresAt time.Time) error {
	cacheKey := r.getCacheKey(role.Role.ID.String())

	redisView, err := converter.CreateEnrichedRoleCacheView(role)
	if err != nil {
		return fmt.Errorf("failed to create cache view: %w", err)
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return fmt.Errorf("expiration time is in the past")
	}

	err = r.redis.HashSet(ctx, cacheKey, redisView)
	if err != nil {
		return fmt.Errorf("failed to store role in cache: %w", err)
	}

	if err = r.redis.Expire(ctx, cacheKey, ttl); err != nil {
		return fmt.Errorf("failed to set expiration: %w", err)
	}

	return nil
}
