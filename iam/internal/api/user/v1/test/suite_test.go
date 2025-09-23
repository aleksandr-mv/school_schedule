package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	api "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/api/user/v1"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/service/mocks"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
)

type APISuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	userService *mocks.UserService
	api         *api.API
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.userService = mocks.NewUserService(s.T())
	s.api = api.NewAPI(s.userService)
}

func (s *APISuite) TearDownTest() {}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
