package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	api "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/api/auth/v1"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/service/mocks"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
)

type APISuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	authService *mocks.AuthServiceInterface
	api         *api.API
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.authService = mocks.NewAuthServiceInterface(s.T())
	s.api = api.NewAPI(s.authService)
}

func (s *APISuite) TearDownTest() {}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
