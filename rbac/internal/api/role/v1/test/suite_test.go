package role_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	api "github.com/aleksandr-mv/school_schedule/rbac/internal/api/role/v1"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/service/mocks"
)

type APISuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	roleService *mocks.RoleServiceInterface
	api         *api.API
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.roleService = mocks.NewRoleServiceInterface(s.T())
	s.api = api.NewAPI(s.roleService)
}

func (s *APISuite) TearDownTest() {}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
