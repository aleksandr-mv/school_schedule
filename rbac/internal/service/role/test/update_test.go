package role_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestUpdateSuccess() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated administrator role"

	updateRole := &model.UpdateRole{
		ID:          roleID.String(),
		Name:        &name,
		Description: &description,
	}

	s.roleRepository.On("Update", mock.Anything, mock.Anything).Return(nil)

	err := s.service.Update(s.ctx, updateRole)

	assert.NoError(s.T(), err)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestUpdateNotFound() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated administrator role"

	updateRole := &model.UpdateRole{
		ID:          roleID.String(),
		Name:        &name,
		Description: &description,
	}

	s.roleRepository.On("Update", mock.Anything, mock.Anything).Return(model.ErrRoleNotFound)

	err := s.service.Update(s.ctx, updateRole)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrRoleNotFound, err)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestUpdateInternalError() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated administrator role"

	updateRole := &model.UpdateRole{
		ID:          roleID.String(),
		Name:        &name,
		Description: &description,
	}

	s.roleRepository.On("Update", mock.Anything, mock.Anything).Return(model.ErrInternal)

	err := s.service.Update(s.ctx, updateRole)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.roleRepository.AssertExpectations(s.T())
}
