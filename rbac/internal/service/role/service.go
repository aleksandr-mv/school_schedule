package role

import (
	"time"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/service"
)

var _ service.RoleServiceInterface = (*RoleService)(nil)

type RoleService struct {
	roleRepo           repository.RoleRepository
	rolePermissionRepo repository.RolePermissionRepository
	enrichedRoleRepo   repository.EnrichedRoleRepository
	enrichedRoleTTL    time.Duration
}

func NewService(
	roleRepo repository.RoleRepository,
	rolePermissionRepo repository.RolePermissionRepository,
	enrichedRoleRepo repository.EnrichedRoleRepository,
	enrichedRoleTTL time.Duration,
) *RoleService {
	return &RoleService{
		roleRepo:           roleRepo,
		rolePermissionRepo: rolePermissionRepo,
		enrichedRoleRepo:   enrichedRoleRepo,
		enrichedRoleTTL:    enrichedRoleTTL,
	}
}
