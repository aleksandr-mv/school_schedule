package user_role_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	repositoryMocks "github.com/aleksandr-mv/school_schedule/rbac/internal/repository/mocks"
	serviceMocks "github.com/aleksandr-mv/school_schedule/rbac/internal/service/mocks"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/service/user_role"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	userRoleRepository *repositoryMocks.UserRoleRepository
	roleRepository     *repositoryMocks.RoleRepository
	roleService        *serviceMocks.RoleServiceInterface

	service *user_role.UserRoleService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.userRoleRepository = repositoryMocks.NewUserRoleRepository(s.T())
	s.roleRepository = repositoryMocks.NewRoleRepository(s.T())

	// Создаем мок для RoleServiceInterface
	s.roleService = serviceMocks.NewRoleServiceInterface(s.T())
	s.service = user_role.NewService(s.userRoleRepository, s.roleService)
}

func (s *ServiceSuite) SetupTest() {
	s.userRoleRepository.ExpectedCalls = nil
	s.roleRepository.ExpectedCalls = nil
	s.roleService.ExpectedCalls = nil
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
