package v1

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"github.com/aleksandr-mv/school_schedule/iam/internal/service"
)

type API struct {
	authv3.UnimplementedAuthorizationServer
	whoAMIService service.WhoAMIService
}

func NewAPI(whoAMIService service.WhoAMIService) *API {
	return &API{
		whoAMIService: whoAMIService,
	}
}
