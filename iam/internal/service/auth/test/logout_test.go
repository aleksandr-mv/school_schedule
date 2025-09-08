package auth_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
)

func (s *ServiceSuite) TestLogoutSuccess() {
	sessionID := uuid.New()

	s.sessionRepository.On("Delete", mock.Anything, sessionID).Return(nil)

	err := s.service.Logout(s.ctx, sessionID)

	assert.NoError(s.T(), err)

	s.sessionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLogoutSessionNotFound() {
	sessionID := uuid.New()

	s.sessionRepository.On("Delete", mock.Anything, sessionID).Return(model.ErrSessionNotFound)

	err := s.service.Logout(s.ctx, sessionID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrFailedToDeleteSession, err)

	s.sessionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLogoutInternalError() {
	sessionID := uuid.New()

	s.sessionRepository.On("Delete", mock.Anything, sessionID).Return(model.ErrInternal)

	err := s.service.Logout(s.ctx, sessionID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrFailedToDeleteSession, err)

	s.sessionRepository.AssertExpectations(s.T())
}
