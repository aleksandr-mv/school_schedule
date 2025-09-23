package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	repositoryMocks "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/mocks"
	serviceMocks "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/service/mocks"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/service/user"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	userRepository         *repositoryMocks.UserRepository
	notificationRepository *repositoryMocks.NotificationRepository
	userProducerService    *serviceMocks.UserProducerService

	service *user.UserService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.userRepository = repositoryMocks.NewUserRepository(s.T())
	s.notificationRepository = repositoryMocks.NewNotificationRepository(s.T())
	s.userProducerService = serviceMocks.NewUserProducerService(s.T())

	s.service = user.NewService(s.userRepository, s.notificationRepository, s.userProducerService)
}

func (s *ServiceSuite) SetupTest() {
	s.userRepository.ExpectedCalls = nil
	s.notificationRepository.ExpectedCalls = nil
	s.userProducerService.ExpectedCalls = nil
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
