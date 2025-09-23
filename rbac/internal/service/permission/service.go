package permission

import (
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service"
)

var _ service.PermissionServiceInterface = (*PermissionService)(nil)

type PermissionService struct {
	permissionRepo repository.PermissionRepository
}

func NewService(permissionRepo repository.PermissionRepository) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
	}
}
