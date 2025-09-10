package user_role

import (
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/service"
)

var _ service.UserRoleServiceInterface = (*UserRoleService)(nil)

type UserRoleService struct {
	userRoleRepo repository.UserRoleRepository
}

func NewService(userRoleRepo repository.UserRoleRepository) *UserRoleService {
	return &UserRoleService{
		userRoleRepo: userRoleRepo,
	}
}
