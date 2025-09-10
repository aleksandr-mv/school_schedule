package permission

import (
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/service"
)

var _ service.PermissionServiceInterface = (*PermissionService)(nil)

type PermissionService struct {
	permissionRepo     repository.PermissionRepository
	rolePermissionRepo repository.RolePermissionRepository
	userRoleRepo       repository.UserRoleRepository
}

func NewService(permissionRepo repository.PermissionRepository,
	rolePermissionRepo repository.RolePermissionRepository,
	userRoleRepo repository.UserRoleRepository) *PermissionService {
	return &PermissionService{
		permissionRepo:     permissionRepo,
		rolePermissionRepo: rolePermissionRepo,
		userRoleRepo:       userRoleRepo,
	}
}
