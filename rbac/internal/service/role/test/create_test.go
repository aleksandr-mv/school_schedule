package role_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestCreateSuccess() {
	roleID := uuid.New()

	s.roleRepository.On("Create", mock.Anything, "admin", "Administrator role").Return(roleID, nil)

	result, err := s.service.Create(s.ctx, "admin", "Administrator role")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), roleID, result)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreateRepositoryError() {
	s.roleRepository.On("Create", mock.Anything, "admin", "Administrator role").Return(uuid.Nil, model.ErrRoleAlreadyExists)

	result, err := s.service.Create(s.ctx, "admin", "Administrator role")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrRoleAlreadyExists, err)
	assert.Equal(s.T(), uuid.Nil, result)

	s.roleRepository.AssertExpectations(s.T())
}
