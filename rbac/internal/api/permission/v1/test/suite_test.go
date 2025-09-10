package permission_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	api "github.com/aleksandr-mv/school_schedule/rbac/internal/api/permission/v1"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/service/mocks"
)

type APISuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	permissionService *mocks.PermissionServiceInterface
	api               *api.API
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.permissionService = mocks.NewPermissionServiceInterface(s.T())
	s.api = api.NewAPI(s.permissionService)
}

func (s *APISuite) TearDownTest() {}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
