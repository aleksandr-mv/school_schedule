package enriched_role

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/converter"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, id string) (*model.EnrichedRole, error) {
	cacheKey := r.getCacheKey(id)

	values, err := r.redis.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, fmt.Errorf("role not found in cache")
		}
		return nil, fmt.Errorf("failed to get role from cache: %w", err)
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("role not found in cache")
	}

	var roleCacheView repoModel.EnrichedRoleCacheView
	err = redigo.ScanStruct(values, &roleCacheView)
	if err != nil {
		return nil, fmt.Errorf("failed to scan cache data: %w", err)
	}

	enrichedRole, err := converter.EnrichedRoleFromRedis(roleCacheView.RoleRedisView)
	if err != nil {
		return nil, fmt.Errorf("failed to convert role data: %w", err)
	}

	return enrichedRole, nil
}
