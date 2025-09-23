package v1

import (
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service"
	rolePermissionV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/role_permission/v1"
)

var _ rolePermissionV1.RolePermissionServiceServer = (*API)(nil)

type API struct {
	rolePermissionV1.UnimplementedRolePermissionServiceServer
	rolePermissionService service.RolePermissionServiceInterface
}

func NewAPI(srv service.RolePermissionServiceInterface) *API {
	return &API{
		rolePermissionService: srv,
	}
}
