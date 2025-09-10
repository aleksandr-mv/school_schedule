package converter

import (
	"strings"

	"github.com/google/uuid"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	repoModel "github.com/aleksandr-mv/school_schedule/rbac/internal/repository/model"
)

// CreateRoleToRepo преобразует доменную модель создания роли в модель репозитория
func CreateRoleToRepo(id uuid.UUID, createRole *model.CreateRole) *repoModel.Role {
	return &repoModel.Role{
		ID:          id,
		Name:        strings.ToLower(createRole.Name),
		Description: createRole.Description,
	}
}

// RoleToDomain преобразует модель репозитория в доменную модель
func RoleToDomain(repoRole *repoModel.Role) *model.Role {
	return &model.Role{
		ID:          repoRole.ID,
		Name:        repoRole.Name,
		Description: repoRole.Description,
		CreatedAt:   repoRole.CreatedAt,
		UpdatedAt:   repoRole.UpdatedAt,
	}
}

// RolesToDomain преобразует массив моделей репозитория в доменные модели
func RolesToDomain(repoRoles []repoModel.Role) []*model.Role {
	result := make([]*model.Role, 0, len(repoRoles))
	for _, role := range repoRoles {
		result = append(result, RoleToDomain(&role))
	}
	return result
}

// UpdateRoleToRepo преобразует параметры обновления роли в модель репозитория
func UpdateRoleToRepo(name *string, description *string) (map[string]interface{}, error) {
	updates := make(map[string]interface{})

	if name != nil {
		updates["name"] = strings.ToLower(*name)
	}

	if description != nil {
		updates["description"] = *description
	}

	return updates, nil
}
