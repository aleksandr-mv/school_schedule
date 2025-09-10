package permission_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestListPermissionsSuccess() {
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

	resource := "users"
	filter := &model.ListPermissionsFilter{
		Resource: &resource,
		Action:   nil,
	}

	s.permissionRepository.On("List", mock.Anything, mock.MatchedBy(func(f *model.ListPermissionsFilter) bool {
		return f.Resource != nil && *f.Resource == "users" && f.Action == nil
	})).Return(expectedPermissions, nil)

	result, err := s.service.ListPermissions(s.ctx, filter)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), expectedPermissions, result)

	s.permissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListPermissionsWithActionFilter() {
	permission1 := &model.Permission{
		ID:       uuid.New(),
		Resource: "users",
		Action:   "read",
	}
	expectedPermissions := []*model.Permission{permission1}

	resource := "users"
	action := "read"
	filter := &model.ListPermissionsFilter{
		Resource: &resource,
		Action:   &action,
	}

	s.permissionRepository.On("List", mock.Anything, mock.MatchedBy(func(f *model.ListPermissionsFilter) bool {
		return f.Resource != nil && *f.Resource == "users" && f.Action != nil && *f.Action == "read"
	})).Return(expectedPermissions, nil)

	result, err := s.service.ListPermissions(s.ctx, filter)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 1)
	assert.Equal(s.T(), expectedPermissions, result)

	s.permissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListPermissionsEmptyResult() {
	resource := "nonexistent"
	filter := &model.ListPermissionsFilter{
		Resource: &resource,
		Action:   nil,
	}

	s.permissionRepository.On("List", mock.Anything, mock.Anything).Return([]*model.Permission{}, nil)

	result, err := s.service.ListPermissions(s.ctx, filter)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 0)

	s.permissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListPermissionsRepositoryError() {
	resource := "users"
	filter := &model.ListPermissionsFilter{
		Resource: &resource,
		Action:   nil,
	}

	s.permissionRepository.On("List", mock.Anything, mock.Anything).Return(nil, model.ErrInternal)

	result, err := s.service.ListPermissions(s.ctx, filter)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)
	assert.Nil(s.T(), result)

	s.permissionRepository.AssertExpectations(s.T())
}
