package permission_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestRevokePermissionFromRoleSuccess() {
	roleID := "role123"
	permissionID := "permission456"

	s.rolePermissionRepository.On("RevokePermissionFromRole", mock.Anything, roleID, permissionID).Return(nil)

	err := s.service.RevokePermissionFromRole(s.ctx, roleID, permissionID)

	assert.NoError(s.T(), err)

	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRevokePermissionFromRoleNotAssigned() {
	roleID := "role123"
	permissionID := "permission456"

	s.rolePermissionRepository.On("RevokePermissionFromRole", mock.Anything, roleID, permissionID).Return(model.ErrPermissionNotAssigned)

	err := s.service.RevokePermissionFromRole(s.ctx, roleID, permissionID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrPermissionNotAssigned, err)

	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRevokePermissionFromRoleRepositoryError() {
	roleID := "role123"
	permissionID := "permission456"

	s.rolePermissionRepository.On("RevokePermissionFromRole", mock.Anything, roleID, permissionID).Return(model.ErrInternal)

	err := s.service.RevokePermissionFromRole(s.ctx, roleID, permissionID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.rolePermissionRepository.AssertExpectations(s.T())
}
