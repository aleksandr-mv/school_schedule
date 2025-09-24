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

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return fmt.Errorf("expiration time is in the past")
	}

	roleConverter := converter.NewEnrichedRoleCacheConverter()
	data, err := roleConverter.ToCache(role)
	if err != nil {
		return fmt.Errorf("failed to convert to protobuf: %w", err)
	}

	err = r.redis.Set(ctx, cacheKey, data, ttl)
	if err != nil {
		return fmt.Errorf("failed to store role in cache: %w", err)
	}

	return nil
}
