package converter

import (
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/model"
)

// PermissionToDomain преобразует модель репозитория в доменную модель
func PermissionToDomain(repoPermission *repoModel.Permission) *model.Permission {
	return &model.Permission{
		ID:       repoPermission.ID,
		Resource: repoPermission.Resource,
		Action:   repoPermission.Action,
	}
}

// PermissionsToDomain преобразует массив моделей репозитория в доменные модели
func PermissionsToDomain(repoPermissions []repoModel.Permission) []*model.Permission {
	result := make([]*model.Permission, 0, len(repoPermissions))
	for _, permission := range repoPermissions {
		result = append(result, PermissionToDomain(&permission))
	}
	return result
}
