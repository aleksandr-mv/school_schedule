package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	api "github.com/aleksandr-mv/school_schedule/iam/internal/api/user/v1"
	"github.com/aleksandr-mv/school_schedule/iam/internal/service/mocks"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

type APISuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	userService *mocks.UserServiceInterface
	api         *api.API
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.userService = mocks.NewUserServiceInterface(s.T())
	s.api = api.NewAPI(s.userService)
}

func (s *APISuite) TearDownTest() {}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
