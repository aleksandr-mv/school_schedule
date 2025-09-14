package model

type EnrichedRoleCacheView struct {
	RoleRedisView EnrichedRoleRedisView `redis:"role"`
}

type EnrichedRoleRedisView struct {
	ID          string `redis:"id"`
	Name        string `redis:"name"`
	Description string `redis:"description"`
	CreatedAt   string `redis:"created_at"`
	UpdatedAt   string `redis:"updated_at"`
	DeletedAt   string `redis:"deleted_at"`
	Permissions string `redis:"permissions"`
}
