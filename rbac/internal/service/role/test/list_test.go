package role_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestListRolesSuccess() {
	role1 := &model.Role{
		ID:          uuid.New(),
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
	}
	role2 := &model.Role{
		ID:          uuid.New(),
		Name:        "user",
		Description: "User role",
		CreatedAt:   time.Now(),
	}
	expectedRoles := []*model.Role{role1, role2}

	s.roleRepository.On("List", mock.Anything, "admin").Return(expectedRoles, nil)

	result, err := s.service.ListRoles(s.ctx, "admin")

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), expectedRoles, result)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListRolesEmptyFilter() {
	role1 := &model.Role{
		ID:          uuid.New(),
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
	}
	role2 := &model.Role{
		ID:          uuid.New(),
		Name:        "user",
		Description: "User role",
		CreatedAt:   time.Now(),
	}
	expectedRoles := []*model.Role{role1, role2}

	s.roleRepository.On("List", mock.Anything, "").Return(expectedRoles, nil)

	result, err := s.service.ListRoles(s.ctx, "")

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), expectedRoles, result)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListRolesEmptyResult() {
	s.roleRepository.On("List", mock.Anything, "nonexistent").Return([]*model.Role{}, nil)

	result, err := s.service.ListRoles(s.ctx, "nonexistent")

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 0)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListRolesRepositoryError() {
	s.roleRepository.On("List", mock.Anything, "admin").Return(nil, model.ErrInternal)

	result, err := s.service.ListRoles(s.ctx, "admin")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)
	assert.Nil(s.T(), result)

	s.roleRepository.AssertExpectations(s.T())
}
