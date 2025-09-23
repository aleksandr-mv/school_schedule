package user_role

import (
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service"
)

var _ service.UserRoleServiceInterface = (*UserRoleService)(nil)

type UserRoleService struct {
	userRoleRepo repository.UserRoleRepository
	roleService  service.RoleServiceInterface
}

func NewService(userRoleRepo repository.UserRoleRepository, roleService service.RoleServiceInterface) *UserRoleService {
	return &UserRoleService{
		userRoleRepo: userRoleRepo,
		roleService:  roleService,
	}
}
