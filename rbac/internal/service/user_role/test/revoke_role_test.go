package user_role_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestRevokeSuccess() {
	userID := "user123"
	roleID := "role456"

	s.userRoleRepository.On("Revoke", mock.Anything, userID, roleID).Return(nil)

	err := s.service.Revoke(s.ctx, userID, roleID)

	assert.NoError(s.T(), err)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRevokeNotAssigned() {
	userID := "user123"
	roleID := "role456"

	s.userRoleRepository.On("Revoke", mock.Anything, userID, roleID).Return(model.ErrRoleNotAssigned)

	err := s.service.Revoke(s.ctx, userID, roleID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrRoleNotAssigned, err)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRevokeRepositoryError() {
	userID := "user123"
	roleID := "role456"

	s.userRoleRepository.On("Revoke", mock.Anything, userID, roleID).Return(model.ErrInternal)

	err := s.service.Revoke(s.ctx, userID, roleID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.userRoleRepository.AssertExpectations(s.T())
}
