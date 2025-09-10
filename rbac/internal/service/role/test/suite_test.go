package role_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository/mocks"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/service/role"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	roleRepository *mocks.RoleRepository

	service *role.RoleService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.roleRepository = mocks.NewRoleRepository(s.T())

	s.service = role.NewService(s.roleRepository)
}

func (s *ServiceSuite) SetupTest() {
	s.roleRepository.ExpectedCalls = nil
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
