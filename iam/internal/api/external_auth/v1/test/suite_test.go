package v1_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	externalAuthV1 "github.com/aleksandr-mv/school_schedule/iam/internal/api/external_auth/v1"
	serviceMocks "github.com/aleksandr-mv/school_schedule/iam/internal/service/mocks"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

type APISuite struct {
	suite.Suite
	ctx context.Context

	api           *externalAuthV1.API
	whoAMIService *serviceMocks.WhoAMIService
}

func (s *APISuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}
}

func (s *APISuite) SetupTest() {
	s.whoAMIService = serviceMocks.NewWhoAMIService(s.T())
	s.api = externalAuthV1.NewAPI(s.whoAMIService)
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}
