package permission_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestGetPermissionByIDSuccess() {
	permissionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	expectedPermission := &model.Permission{
		ID:       permissionID,
		Resource: "users",
		Action:   "read",
	}

	s.permissionRepository.On("Get", mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(expectedPermission, nil)

	result, err := s.service.GetPermission(s.ctx, "123e4567-e89b-12d3-a456-426614174000")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedPermission.ID, result.ID)
	assert.Equal(s.T(), expectedPermission.Resource, result.Resource)
	assert.Equal(s.T(), expectedPermission.Action, result.Action)

	s.permissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetPermissionByNameSuccess() {
	permissionID := uuid.New()
	expectedPermission := &model.Permission{
		ID:       permissionID,
		Resource: "users",
		Action:   "read",
	}

	s.permissionRepository.On("Get", mock.Anything, "users:read").Return(expectedPermission, nil)

	result, err := s.service.GetPermission(s.ctx, "users:read")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedPermission.ID, result.ID)
	assert.Equal(s.T(), expectedPermission.Resource, result.Resource)
	assert.Equal(s.T(), expectedPermission.Action, result.Action)

	s.permissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetPermissionNotFound() {
	s.permissionRepository.On("Get", mock.Anything, "nonexistent").Return(nil, model.ErrPermissionNotFound)

	result, err := s.service.GetPermission(s.ctx, "nonexistent")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrPermissionNotFound, err)
	assert.Nil(s.T(), result)

	s.permissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetPermissionRepositoryError() {
	s.permissionRepository.On("Get", mock.Anything, "users:read").Return(nil, model.ErrInternal)

	result, err := s.service.GetPermission(s.ctx, "users:read")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)
	assert.Nil(s.T(), result)

	s.permissionRepository.AssertExpectations(s.T())
}
