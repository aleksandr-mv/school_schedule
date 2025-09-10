package v1

import (
	"github.com/aleksandr-mv/school_schedule/rbac/internal/service"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
)

// API реализует RoleService gRPC сервер
type API struct {
	roleV1.UnimplementedRoleServiceServer
	roleService service.RoleServiceInterface
}

// NewAPI создает новый экземпляр API для RoleService
func NewAPI(roleService service.RoleServiceInterface) *API {
	return &API{
		roleService: roleService,
	}
}
