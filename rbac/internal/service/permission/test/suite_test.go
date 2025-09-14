package permission_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository/mocks"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/service/permission"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	permissionRepository *mocks.PermissionRepository

	service *permission.PermissionService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.permissionRepository = mocks.NewPermissionRepository(s.T())

	s.service = permission.NewService(s.permissionRepository)
}

func (s *ServiceSuite) SetupTest() {
	s.permissionRepository.ExpectedCalls = nil
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
