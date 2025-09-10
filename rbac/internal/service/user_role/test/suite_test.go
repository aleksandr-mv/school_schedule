package user_role_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository/mocks"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/service/user_role"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	userRoleRepository *mocks.UserRoleRepository

	service *user_role.UserRoleService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.userRoleRepository = mocks.NewUserRoleRepository(s.T())

	s.service = user_role.NewService(s.userRoleRepository)
}

func (s *ServiceSuite) SetupTest() {
	s.userRoleRepository.ExpectedCalls = nil
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
