package user_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
)

func (s *ServiceSuite) TestGetUserSuccess() {
	userID := uuid.New()

	user := &model.User{
		ID:        userID,
		Login:     "testuser123",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	notificationMethods := []*model.NotificationMethod{
		{
			UserID:       userID,
			ProviderName: "email",
			Target:       "test@example.com",
			CreatedAt:    time.Now(),
			UpdatedAt:    nil,
		},
	}

	userWithNotifications := *user
	userWithNotifications.NotificationMethods = notificationMethods

	s.userRepository.On("Get", mock.Anything, userID.String()).Return(user, nil)
	s.notificationRepository.On("GetByUser", mock.Anything, userID).Return(notificationMethods, nil)

	result, err := s.service.GetUser(s.ctx, userID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), userID, result.ID)
	assert.Equal(s.T(), user.Login, result.Login)
	assert.Equal(s.T(), user.Email, result.Email)
	assert.Equal(s.T(), notificationMethods, result.NotificationMethods)

	s.userRepository.AssertExpectations(s.T())
	s.notificationRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetUserNotFound() {
	userID := uuid.New()

	s.userRepository.On("Get", mock.Anything, userID.String()).Return(nil, model.ErrUserNotFound)

	result, err := s.service.GetUser(s.ctx, userID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrUserNotFound, err)
	assert.Nil(s.T(), result)

	s.userRepository.AssertExpectations(s.T())
	s.notificationRepository.AssertNotCalled(s.T(), "GetByUser")
}

func (s *ServiceSuite) TestGetUserNotificationError() {
	userID := uuid.New()

	user := &model.User{
		ID:        userID,
		Login:     "testuser123",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	s.userRepository.On("Get", mock.Anything, userID.String()).Return(user, nil)
	s.notificationRepository.On("GetByUser", mock.Anything, userID).Return(nil, model.ErrFailedToListNotifications)

	result, err := s.service.GetUser(s.ctx, userID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrFailedToListNotifications, err)
	assert.Nil(s.T(), result)

	s.userRepository.AssertExpectations(s.T())
	s.notificationRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetUserRepositoryError() {
	userID := uuid.New()

	s.userRepository.On("Get", mock.Anything, userID.String()).Return(nil, model.ErrInternal)

	result, err := s.service.GetUser(s.ctx, userID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)
	assert.Nil(s.T(), result)

	s.userRepository.AssertExpectations(s.T())
	s.notificationRepository.AssertNotCalled(s.T(), "GetByUser")
}
