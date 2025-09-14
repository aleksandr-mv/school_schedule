package enriched_role

import "fmt"

const enrichedRoleCachePrefix = "enriched_role"

func (r *repository) getCacheKey(roleID string) string {
	return fmt.Sprintf("%s:%s", enrichedRoleCachePrefix, roleID)
}
