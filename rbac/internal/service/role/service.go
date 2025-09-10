package role

import (
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/service"
)

var _ service.RoleServiceInterface = (*RoleService)(nil)

type RoleService struct {
	roleRepo repository.RoleRepository
}

func NewService(roleRepo repository.RoleRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
	}
}
