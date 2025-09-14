package permission_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestListSuccess() {
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

	s.permissionRepository.On("List", mock.Anything).Return(expectedPermissions, nil).Once()

	permissions, err := s.service.List(s.ctx)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), permissions)
	assert.Len(s.T(), permissions, 2)

	assert.Equal(s.T(), permission1.ID, permissions[0].ID)
	assert.Equal(s.T(), permission1.Resource, permissions[0].Resource)
	assert.Equal(s.T(), permission1.Action, permissions[0].Action)

	assert.Equal(s.T(), permission2.ID, permissions[1].ID)
	assert.Equal(s.T(), permission2.Resource, permissions[1].Resource)
	assert.Equal(s.T(), permission2.Action, permissions[1].Action)

	s.permissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListEmptyResult() {
	expectedPermissions := []*model.Permission{}

	s.permissionRepository.On("List", mock.Anything).Return(expectedPermissions, nil).Once()

	permissions, err := s.service.List(s.ctx)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), permissions)
	assert.Len(s.T(), permissions, 0)

	s.permissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListInternalError() {
	s.permissionRepository.On("List", mock.Anything).Return(nil, model.ErrInternal).Once()

	permissions, err := s.service.List(s.ctx)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), permissions)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.permissionRepository.AssertExpectations(s.T())
}
