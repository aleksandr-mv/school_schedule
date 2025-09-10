package permission_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestListRolesByPermissionSuccess() {
	role1 := &model.Role{
		ID:        uuid.New(),
		Name:      "admin",
		CreatedAt: time.Now(),
	}
	role2 := &model.Role{
		ID:        uuid.New(),
		Name:      "moderator",
		CreatedAt: time.Now(),
	}
	expectedRoles := []*model.Role{role1, role2}

	permissionValue := "users:read"

	s.rolePermissionRepository.On("ListRolesByPermission", mock.Anything, permissionValue).Return(expectedRoles, nil)

	result, err := s.service.ListRolesByPermission(s.ctx, permissionValue)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), expectedRoles, result)

	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListRolesByPermissionEmptyResult() {
	permissionValue := "nonexistent:action"

	s.rolePermissionRepository.On("ListRolesByPermission", mock.Anything, permissionValue).Return([]*model.Role{}, nil)

	result, err := s.service.ListRolesByPermission(s.ctx, permissionValue)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 0)

	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListRolesByPermissionRepositoryError() {
	permissionValue := "users:read"

	s.rolePermissionRepository.On("ListRolesByPermission", mock.Anything, permissionValue).Return(nil, model.ErrInternal)

	result, err := s.service.ListRolesByPermission(s.ctx, permissionValue)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)
	assert.Nil(s.T(), result)

	s.rolePermissionRepository.AssertExpectations(s.T())
}
