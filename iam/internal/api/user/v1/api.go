package v1

import (
	"github.com/aleksandr-mv/school_schedule/iam/internal/service"
	userV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user/v1"
)

type API struct {
	userV1.UnimplementedUserServiceServer
	userService service.UserServiceInterface
}

func NewAPI(userService service.UserServiceInterface) *API {
	return &API{
		userService: userService,
	}
}
