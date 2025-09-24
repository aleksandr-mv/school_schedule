package enriched_role

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, id string) (*model.EnrichedRole, error) {
	cacheKey := r.getCacheKey(id)

	data, err := r.redis.Get(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("role not found in cache")
		}
		return nil, fmt.Errorf("failed to get role from cache: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("role not found in cache")
	}

	roleConverter := converter.NewEnrichedRoleCacheConverter()
	enrichedRole, err := roleConverter.FromCache(data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert from protobuf: %w", err)
	}

	return enrichedRole, nil
}
