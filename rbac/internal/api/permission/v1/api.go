package v1

import (
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service"
	permissionV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/permission/v1"
)

var _ permissionV1.PermissionServiceServer = (*API)(nil)

type API struct {
	permissionV1.UnimplementedPermissionServiceServer
	permissionService service.PermissionServiceInterface
}

func NewAPI(permissionService service.PermissionServiceInterface) *API {
	return &API{
		permissionService: permissionService,
	}
}
