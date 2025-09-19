package whoami

import (
	"github.com/aleksandr-mv/school_schedule/iam/internal/client/grpc"
	"github.com/aleksandr-mv/school_schedule/iam/internal/repository"
	def "github.com/aleksandr-mv/school_schedule/iam/internal/service"
)

var _ def.WhoAMIService = (*WhoAMIService)(nil)

type WhoAMIService struct {
	sessionRepository repository.SessionRepository
	rbacClient        grpc.RBACClient
}

func NewService(
	sessionRepository repository.SessionRepository,
	rbacClient grpc.RBACClient,
) *WhoAMIService {
	return &WhoAMIService{
		sessionRepository: sessionRepository,
		rbacClient:        rbacClient,
	}
}
