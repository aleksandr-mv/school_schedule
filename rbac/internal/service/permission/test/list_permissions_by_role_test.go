package permission_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestListPermissionsByRoleSuccess() {
	permission1 := &model.Permission{
		ID:       uuid.New(),
		Resource: "users",
		Action:   "read",
	}
	permission2 := &model.Permission{
		ID:       uuid.New(),
		Resource: "users",
		Action:   "write",
	}
	expectedPermissions := []*model.Permission{permission1, permission2}

	roleValue := "admin"

	s.rolePermissionRepository.On("ListPermissionsByRole", mock.Anything, roleValue).Return(expectedPermissions, nil)

	result, err := s.service.ListPermissionsByRole(s.ctx, roleValue)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), expectedPermissions, result)

	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListPermissionsByRoleEmptyResult() {
	roleValue := "user"

	s.rolePermissionRepository.On("ListPermissionsByRole", mock.Anything, roleValue).Return([]*model.Permission{}, nil)

	result, err := s.service.ListPermissionsByRole(s.ctx, roleValue)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 0)

	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListPermissionsByRoleRepositoryError() {
	roleValue := "admin"

	s.rolePermissionRepository.On("ListPermissionsByRole", mock.Anything, roleValue).Return(nil, model.ErrInternal)

	result, err := s.service.ListPermissionsByRole(s.ctx, roleValue)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)
	assert.Nil(s.T(), result)

	s.rolePermissionRepository.AssertExpectations(s.T())
}
