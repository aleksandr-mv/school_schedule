package role_permission

import (
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service"
)

var _ service.RolePermissionServiceInterface = (*RolePermissionService)(nil)

type RolePermissionService struct {
	rolePermissionRepo repository.RolePermissionRepository
}

func NewService(rolePermissionRepo repository.RolePermissionRepository) *RolePermissionService {
	return &RolePermissionService{
		rolePermissionRepo: rolePermissionRepo,
	}
}
