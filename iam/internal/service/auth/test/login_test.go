package auth_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
)

// Создаем валидный bcrypt хеш для тестов
var validPasswordHash = func() string {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123456"), bcrypt.MinCost)
	return string(hash)
}()

func (s *ServiceSuite) TestLoginSuccess() {
	userID := uuid.New()
	sessionID := uuid.New()

	credentials := &model.LoginCredentials{
		Login:    "testuser123",
		Password: "password123456",
	}

	user := &model.User{
		ID:           userID,
		Login:        "testuser123",
		Email:        "test@example.com",
		PasswordHash: validPasswordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    nil,
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

	s.userRepository.On("Get", mock.Anything, credentials.Login).Return(user, nil)
	s.notificationRepository.On("GetByUser", mock.Anything, userID).Return(notificationMethods, nil)
	s.sessionRepository.On("Create", mock.Anything, userWithNotifications, mock.AnythingOfType("time.Time")).Return(sessionID, nil)

	result, err := s.service.Login(s.ctx, credentials)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), sessionID, result)

	s.userRepository.AssertExpectations(s.T())
	s.notificationRepository.AssertExpectations(s.T())
	s.sessionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLoginUserNotFound() {
	credentials := &model.LoginCredentials{
		Login:    "nonexistent",
		Password: "password123456",
	}

	s.userRepository.On("Get", mock.Anything, credentials.Login).Return(nil, model.ErrUserNotFound)

	result, err := s.service.Login(s.ctx, credentials)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInvalidCredentials, err)
	assert.Equal(s.T(), uuid.Nil, result)

	s.userRepository.AssertExpectations(s.T())
	s.notificationRepository.AssertNotCalled(s.T(), "GetByUser")
	s.sessionRepository.AssertNotCalled(s.T(), "Create")
}

func (s *ServiceSuite) TestLoginInvalidPassword() {
	userID := uuid.New()

	credentials := &model.LoginCredentials{
		Login:    "testuser123",
		Password: "wrongpassword",
	}

	user := &model.User{
		ID:           userID,
		Login:        "testuser123",
		Email:        "test@example.com",
		PasswordHash: validPasswordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    nil,
	}

	s.userRepository.On("Get", mock.Anything, credentials.Login).Return(user, nil)

	result, err := s.service.Login(s.ctx, credentials)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInvalidCredentials, err)
	assert.Equal(s.T(), uuid.Nil, result)

	s.userRepository.AssertExpectations(s.T())
	s.notificationRepository.AssertNotCalled(s.T(), "GetByUser")
	s.sessionRepository.AssertNotCalled(s.T(), "Create")
}

func (s *ServiceSuite) TestLoginNotificationError() {
	userID := uuid.New()

	credentials := &model.LoginCredentials{
		Login:    "testuser123",
		Password: "password123456",
	}

	user := &model.User{
		ID:           userID,
		Login:        "testuser123",
		Email:        "test@example.com",
		PasswordHash: validPasswordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    nil,
	}

	s.userRepository.On("Get", mock.Anything, credentials.Login).Return(user, nil)
	s.notificationRepository.On("GetByUser", mock.Anything, userID).Return(nil, model.ErrFailedToListNotifications)

	result, err := s.service.Login(s.ctx, credentials)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrFailedToListNotifications, err)
	assert.Equal(s.T(), uuid.Nil, result)

	s.userRepository.AssertExpectations(s.T())
	s.notificationRepository.AssertExpectations(s.T())
	s.sessionRepository.AssertNotCalled(s.T(), "Create")
}

func (s *ServiceSuite) TestLoginSessionCreationError() {
	userID := uuid.New()

	credentials := &model.LoginCredentials{
		Login:    "testuser123",
		Password: "password123456",
	}

	user := &model.User{
		ID:           userID,
		Login:        "testuser123",
		Email:        "test@example.com",
		PasswordHash: validPasswordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    nil,
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

	s.userRepository.On("Get", mock.Anything, credentials.Login).Return(user, nil)
	s.notificationRepository.On("GetByUser", mock.Anything, userID).Return(notificationMethods, nil)
	s.sessionRepository.On("Create", mock.Anything, userWithNotifications, mock.AnythingOfType("time.Time")).Return(uuid.Nil, model.ErrFailedToCreateSession)

	result, err := s.service.Login(s.ctx, credentials)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrFailedToCreateSession, err)
	assert.Equal(s.T(), uuid.Nil, result)

	s.userRepository.AssertExpectations(s.T())
	s.notificationRepository.AssertExpectations(s.T())
	s.sessionRepository.AssertExpectations(s.T())
}
