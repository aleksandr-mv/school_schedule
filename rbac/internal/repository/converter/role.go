package converter

import (
	"strings"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/model"
)

// RoleToDomain преобразует модель репозитория в доменную модель
func RoleToDomain(repoRole *repoModel.Role) *model.Role {
	return &model.Role{
		ID:          repoRole.ID,
		Name:        repoRole.Name,
		Description: repoRole.Description,
		CreatedAt:   repoRole.CreatedAt,
		UpdatedAt:   repoRole.UpdatedAt,
		DeletedAt:   repoRole.DeletedAt,
	}
}

// UpdateRoleToRepo преобразует параметры обновления роли в модель репозитория
func UpdateRoleToRepo(name, description *string) (map[string]interface{}, error) {
	updates := make(map[string]interface{})

	if name != nil {
		updates["name"] = strings.ToLower(*name)
	}

	if description != nil {
		updates["description"] = *description
	}

	return updates, nil
}

// RolesToDomain преобразует список моделей репозитория в список доменных моделей
func RolesToDomain(rows []repoModel.Role) ([]*model.Role, error) {
	result := make([]*model.Role, 0, len(rows))
	for _, row := range rows {
		role := RoleToDomain(&row)
		result = append(result, role)
	}
	return result, nil
}
