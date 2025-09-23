package converter

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/model"
)

func EnrichedRoleToRedis(enrichedRole *model.EnrichedRole) (repoModel.EnrichedRoleRedisView, error) {
	permissionsJSON, err := json.Marshal(enrichedRole.Permissions)
	if err != nil {
		return repoModel.EnrichedRoleRedisView{}, err
	}

	roleView := repoModel.EnrichedRoleRedisView{
		ID:          enrichedRole.Role.ID.String(),
		Name:        enrichedRole.Role.Name,
		Description: enrichedRole.Role.Description,
		CreatedAt:   enrichedRole.Role.CreatedAt.Format(time.RFC3339),
		Permissions: string(permissionsJSON),
	}

	if enrichedRole.Role.UpdatedAt != nil {
		roleView.UpdatedAt = enrichedRole.Role.UpdatedAt.Format(time.RFC3339)
	}

	if enrichedRole.Role.DeletedAt != nil {
		roleView.DeletedAt = enrichedRole.Role.DeletedAt.Format(time.RFC3339)
	}

	return roleView, nil
}

func EnrichedRoleFromRedis(redisView repoModel.EnrichedRoleRedisView) (*model.EnrichedRole, error) {
	roleID, err := uuid.Parse(redisView.ID)
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, redisView.CreatedAt)
	if err != nil {
		return nil, err
	}

	var updatedAt *time.Time
	if redisView.UpdatedAt != "" {
		parsed, err := time.Parse(time.RFC3339, redisView.UpdatedAt)
		if err != nil {
			return nil, err
		}
		updatedAt = &parsed
	}

	var deletedAt *time.Time
	if redisView.DeletedAt != "" {
		parsed, err := time.Parse(time.RFC3339, redisView.DeletedAt)
		if err != nil {
			return nil, err
		}
		deletedAt = &parsed
	}

	var permissions []*model.Permission
	if err := json.Unmarshal([]byte(redisView.Permissions), &permissions); err != nil {
		return nil, err
	}

	return &model.EnrichedRole{
		Role: model.Role{
			ID:          roleID,
			Name:        redisView.Name,
			Description: redisView.Description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			DeletedAt:   deletedAt,
		},
		Permissions: permissions,
	}, nil
}

func CreateEnrichedRoleCacheView(enrichedRole *model.EnrichedRole) (repoModel.EnrichedRoleCacheView, error) {
	roleView, err := EnrichedRoleToRedis(enrichedRole)
	if err != nil {
		return repoModel.EnrichedRoleCacheView{}, err
	}

	return repoModel.EnrichedRoleCacheView{
		RoleRedisView: roleView,
	}, nil
}
