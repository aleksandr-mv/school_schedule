package role_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestGetRoleByIDSuccess() {
	roleID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	expectedRole := &model.Role{
		ID:          roleID,
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
	}

	s.roleRepository.On("Get", mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(expectedRole, nil)

	result, err := s.service.GetRole(s.ctx, "123e4567-e89b-12d3-a456-426614174000")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedRole.ID, result.ID)
	assert.Equal(s.T(), expectedRole.Name, result.Name)
	assert.Equal(s.T(), expectedRole.Description, result.Description)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetRoleByNameSuccess() {
	roleID := uuid.New()
	expectedRole := &model.Role{
		ID:          roleID,
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
	}

	s.roleRepository.On("Get", mock.Anything, "admin").Return(expectedRole, nil)

	result, err := s.service.GetRole(s.ctx, "admin")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedRole.ID, result.ID)
	assert.Equal(s.T(), expectedRole.Name, result.Name)
	assert.Equal(s.T(), expectedRole.Description, result.Description)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetRoleNotFound() {
	s.roleRepository.On("Get", mock.Anything, "nonexistent").Return(nil, model.ErrRoleNotFound)

	result, err := s.service.GetRole(s.ctx, "nonexistent")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrRoleNotFound, err)
	assert.Nil(s.T(), result)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetRoleRepositoryError() {
	s.roleRepository.On("Get", mock.Anything, "admin").Return(nil, model.ErrInternal)

	result, err := s.service.GetRole(s.ctx, "admin")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)
	assert.Nil(s.T(), result)

	s.roleRepository.AssertExpectations(s.T())
}
