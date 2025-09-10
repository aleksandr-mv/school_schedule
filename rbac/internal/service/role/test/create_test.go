package role_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestCreateRoleSuccess() {
	roleID := uuid.New()

	createRole := &model.CreateRole{
		Name:        "admin",
		Description: "Administrator role",
	}

	expectedRole := &model.Role{
		ID:          roleID,
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
	}

	s.roleRepository.On("Create", mock.Anything, mock.MatchedBy(func(cr *model.CreateRole) bool {
		return cr.Name == "admin" && cr.Description == "Administrator role"
	})).Return(expectedRole, nil)

	result, err := s.service.CreateRole(s.ctx, createRole)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), roleID, result)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreateRoleValidationError() {
	createRole := &model.CreateRole{
		Name:        "", // Invalid: empty name
		Description: "Empty name role",
	}

	result, err := s.service.CreateRole(s.ctx, createRole)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInvalidCredentials, err)
	assert.Equal(s.T(), uuid.Nil, result)

	s.roleRepository.AssertNotCalled(s.T(), "Create")
}

func (s *ServiceSuite) TestCreateRoleRepositoryError() {
	createRole := &model.CreateRole{
		Name:        "admin",
		Description: "Administrator role",
	}

	s.roleRepository.On("Create", mock.Anything, mock.Anything).Return(nil, model.ErrRoleAlreadyExists)

	result, err := s.service.CreateRole(s.ctx, createRole)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrRoleAlreadyExists, err)
	assert.Equal(s.T(), uuid.Nil, result)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreateRoleNameTooShort() {
	createRole := &model.CreateRole{
		Name:        "a", // Invalid: too short
		Description: "Short name role",
	}

	result, err := s.service.CreateRole(s.ctx, createRole)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInvalidCredentials, err)
	assert.Equal(s.T(), uuid.Nil, result)

	s.roleRepository.AssertNotCalled(s.T(), "Create")
}

func (s *ServiceSuite) TestCreateRoleNameTooLong() {
	createRole := &model.CreateRole{
		Name:        "very_long_role_name_that_exceeds_fifty_characters_limit", // Invalid: too long
		Description: "Long name role",
	}

	result, err := s.service.CreateRole(s.ctx, createRole)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInvalidCredentials, err)
	assert.Equal(s.T(), uuid.Nil, result)

	s.roleRepository.AssertNotCalled(s.T(), "Create")
}
