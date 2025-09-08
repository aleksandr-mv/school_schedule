package session

import "fmt"

const (
	cacheKeyPrefix = "session:"
)

func (r *sessionRepository) getCacheKey(id string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, id)
}
