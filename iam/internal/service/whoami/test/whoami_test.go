package whoami_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
)

func (s *ServiceSuite) TestWhoamiSuccess() {
	userID := uuid.New()
	sessionID := uuid.New()
	expiresAt := time.Now().Add(time.Hour)

	expectedWhoami := &model.WhoAMI{
		Session: model.Session{
			ID:        sessionID,
			ExpiresAt: expiresAt,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		User: model.User{
			ID:        userID,
			Login:     "testuser",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: nil,
			NotificationMethods: []*model.NotificationMethod{
				{
					UserID:       userID,
					ProviderName: "email",
					Target:       "test@example.com",
					CreatedAt:    time.Now(),
					UpdatedAt:    nil,
				},
			},
		},
	}

	roleID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	permissionID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")

	expectedWhoami.RolesWithPermissions = []*model.RoleWithPermissions{
		{
			Role: &model.Role{
				ID:   roleID,
				Name: "admin",
			},
			Permissions: []*model.Permission{
				{ID: permissionID, Resource: "user", Action: "read"},
			},
		},
	}

	s.sessionRepository.On("Get", mock.Anything, sessionID).Return(expectedWhoami, nil)

	result, err := s.service.Whoami(s.ctx, sessionID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedWhoami.Session.ID, result.Session.ID)
	assert.Equal(s.T(), expectedWhoami.User.ID, result.User.ID)
	assert.Equal(s.T(), expectedWhoami.User.Login, result.User.Login)
	assert.Equal(s.T(), expectedWhoami.User.Email, result.User.Email)
	assert.Equal(s.T(), expectedWhoami.RolesWithPermissions, result.RolesWithPermissions)

	s.sessionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestWhoamiSessionNotFound() {
	sessionID := uuid.New()

	s.sessionRepository.On("Get", mock.Anything, sessionID).Return(nil, model.ErrSessionNotFound)

	result, err := s.service.Whoami(s.ctx, sessionID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrSessionNotFound, err)
	assert.Nil(s.T(), result)

	s.sessionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestWhoamiSessionExpired() {
	sessionID := uuid.New()
	expiresAt := time.Now().Add(-time.Hour) // Истекшая сессия

	expiredSession := &model.WhoAMI{
		Session: model.Session{
			ID:        sessionID,
			ExpiresAt: expiresAt,
		},
	}

	s.sessionRepository.On("Get", mock.Anything, sessionID).Return(expiredSession, nil)

	result, err := s.service.Whoami(s.ctx, sessionID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrSessionExpired, err)
	assert.Nil(s.T(), result)

	s.sessionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestWhoamiWithoutRoles() {
	userID := uuid.New()
	sessionID := uuid.New()
	expiresAt := time.Now().Add(time.Hour)

	expectedWhoami := &model.WhoAMI{
		Session: model.Session{
			ID:        sessionID,
			ExpiresAt: expiresAt,
		},
		User: model.User{
			ID: userID,
		},
		RolesWithPermissions: []*model.RoleWithPermissions{}, // Пустые роли из Redis
	}

	s.sessionRepository.On("Get", mock.Anything, sessionID).Return(expectedWhoami, nil)

	result, err := s.service.Whoami(s.ctx, sessionID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.RolesWithPermissions)

	s.sessionRepository.AssertExpectations(s.T())
}
