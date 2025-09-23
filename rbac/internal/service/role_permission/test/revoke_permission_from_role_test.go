package role_permission_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestRevokeSuccess() {
	roleID := "role123"
	permissionID := "permission456"

	s.rolePermissionRepository.On("Revoke", mock.Anything, roleID, permissionID).Return(nil)

	err := s.service.Revoke(s.ctx, roleID, permissionID)

	assert.NoError(s.T(), err)

	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRevokeNotAssigned() {
	roleID := "role123"
	permissionID := "permission456"

	s.rolePermissionRepository.On("Revoke", mock.Anything, roleID, permissionID).Return(model.ErrPermissionNotAssigned)

	err := s.service.Revoke(s.ctx, roleID, permissionID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrPermissionNotAssigned, err)

	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRevokeRepositoryError() {
	roleID := "role123"
	permissionID := "permission456"

	s.rolePermissionRepository.On("Revoke", mock.Anything, roleID, permissionID).Return(model.ErrInternal)

	err := s.service.Revoke(s.ctx, roleID, permissionID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.rolePermissionRepository.AssertExpectations(s.T())
}
