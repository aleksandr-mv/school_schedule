package enriched_role

import (
	"github.com/aleksandr-mv/school_schedule/platform/pkg/cache"
	def "github.com/aleksandr-mv/school_schedule/rbac/internal/repository"
)

var _ def.EnrichedRoleRepository = (*repository)(nil)

type repository struct {
	redis cache.RedisClient
}

func NewRepository(redis cache.RedisClient) *repository {
	return &repository{
		redis: redis,
	}
}
