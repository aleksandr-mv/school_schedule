package v1

import (
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/service"
	userV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/user/v1"
)

type API struct {
	userV1.UnimplementedUserServiceServer
	userService service.UserService
}

func NewAPI(userService service.UserService) *API {
	return &API{
		userService: userService,
	}
}
