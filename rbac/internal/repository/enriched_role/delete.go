package enriched_role

import (
	"context"
	"fmt"
)

func (r *repository) Delete(ctx context.Context, id string) error {
	cacheKey := r.getCacheKey(id)
	if err := r.redis.Del(ctx, cacheKey); err != nil {
		return fmt.Errorf("failed to delete role from cache: %w", err)
	}

	return nil
}
