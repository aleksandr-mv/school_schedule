package role_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestUpdateRoleSuccess() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated administrator role"

	updateRole := &model.UpdateRole{
		ID:          roleID.String(),
		Name:        &name,
		Description: &description,
	}

	s.roleRepository.On("Update", mock.Anything, mock.Anything).Return(nil)

	err := s.service.UpdateRole(s.ctx, updateRole)

	assert.NoError(s.T(), err)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestUpdateRoleValidationError() {
	roleID := uuid.New()
	emptyName := ""
	description := "Updated administrator role"

	updateRole := &model.UpdateRole{
		ID:          roleID.String(),
		Name:        &emptyName, // Invalid: empty name
		Description: &description,
	}

	err := s.service.UpdateRole(s.ctx, updateRole)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInvalidCredentials, err)

	s.roleRepository.AssertNotCalled(s.T(), "Update")
}

func (s *ServiceSuite) TestUpdateRoleNotFound() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated administrator role"

	updateRole := &model.UpdateRole{
		ID:          roleID.String(),
		Name:        &name,
		Description: &description,
	}

	s.roleRepository.On("Update", mock.Anything, mock.Anything).Return(model.ErrRoleNotFound)

	err := s.service.UpdateRole(s.ctx, updateRole)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrRoleNotFound, err)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestUpdateRoleInternalError() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated administrator role"

	updateRole := &model.UpdateRole{
		ID:          roleID.String(),
		Name:        &name,
		Description: &description,
	}

	s.roleRepository.On("Update", mock.Anything, mock.Anything).Return(model.ErrInternal)

	err := s.service.UpdateRole(s.ctx, updateRole)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.roleRepository.AssertExpectations(s.T())
}
