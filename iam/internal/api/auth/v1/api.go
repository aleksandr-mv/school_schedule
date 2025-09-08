package v1

import (
	"github.com/aleksandr-mv/school_schedule/iam/internal/service"
	authV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/auth/v1"
)

// API реализует AuthService gRPC сервер
type API struct {
	authV1.UnimplementedAuthServiceServer
	authService service.AuthServiceInterface
}

// NewAPI создает новый экземпляр API для AuthService
func NewAPI(authService service.AuthServiceInterface) *API {
	return &API{
		authService: authService,
	}
}
