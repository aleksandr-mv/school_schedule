package whoami_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	grpcMocks "github.com/aleksandr-mv/school_schedule/iam/internal/client/grpc/mocks"
	repositoryMocks "github.com/aleksandr-mv/school_schedule/iam/internal/repository/mocks"
	"github.com/aleksandr-mv/school_schedule/iam/internal/service/whoami"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context

	service *whoami.WhoAMIService

	sessionRepository *repositoryMocks.SessionRepository
	rbacClient        *grpcMocks.RBACClient
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.sessionRepository = repositoryMocks.NewSessionRepository(s.T())
	s.rbacClient = grpcMocks.NewRBACClient(s.T())

	s.service = whoami.NewService(
		s.sessionRepository,
		s.rbacClient,
	)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
