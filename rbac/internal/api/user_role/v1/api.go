package v1

import (
	"github.com/aleksandr-mv/school_schedule/rbac/internal/service"
	userRoleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user_role/v1"
)

var _ userRoleV1.UserRoleServiceServer = (*API)(nil)

type API struct {
	userRoleV1.UnimplementedUserRoleServiceServer
	userRoleService service.UserRoleServiceInterface
}

func NewAPI(userRoleService service.UserRoleServiceInterface) *API {
	return &API{
		userRoleService: userRoleService,
	}
}
