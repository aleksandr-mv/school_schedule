package whoami_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	repositoryMocks "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/mocks"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/service/whoami"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context //nolint:containedctx // test context setup

	service *whoami.WhoAMIService

	sessionRepository *repositoryMocks.SessionRepository
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.sessionRepository = repositoryMocks.NewSessionRepository(s.T())

	s.service = whoami.NewService(s.sessionRepository)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
