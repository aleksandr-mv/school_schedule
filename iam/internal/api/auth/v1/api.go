package v1

import (
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/service"
	authV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/auth/v1"
)

// API реализует AuthService gRPC сервер
type API struct {
	authV1.UnimplementedAuthServiceServer
	authService   service.AuthService
	whoAMIService service.WhoAMIService
}

// NewAPI создает новый экземпляр API для AuthService
func NewAPI(authService service.AuthService, whoAMIService service.WhoAMIService) *API {
	return &API{
		authService:   authService,
		whoAMIService: whoAMIService,
	}
}
