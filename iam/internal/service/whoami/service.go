package whoami

import (
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository"
	def "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/service"
)

var _ def.WhoAMIService = (*WhoAMIService)(nil)

type WhoAMIService struct {
	sessionRepository repository.SessionRepository
}

func NewService(sessionRepository repository.SessionRepository) *WhoAMIService {
	return &WhoAMIService{
		sessionRepository: sessionRepository,
	}
}
