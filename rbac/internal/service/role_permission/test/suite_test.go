package role_permission_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/mocks"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service/role_permission"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	rolePermissionRepository *mocks.RolePermissionRepository

	service *role_permission.RolePermissionService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.rolePermissionRepository = mocks.NewRolePermissionRepository(s.T())

	s.service = role_permission.NewService(s.rolePermissionRepository)
}

func (s *ServiceSuite) SetupTest() {
	s.rolePermissionRepository.ExpectedCalls = nil
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
