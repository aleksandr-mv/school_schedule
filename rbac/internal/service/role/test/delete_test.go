package role_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestDeleteSuccess() {
	roleID := "123e4567-e89b-12d3-a456-426614174000"

	s.roleRepository.On("Delete", mock.Anything, roleID).Return(nil)

	err := s.service.Delete(s.ctx, roleID)

	assert.NoError(s.T(), err)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestDeleteNotFound() {
	roleID := "123e4567-e89b-12d3-a456-426614174000"

	s.roleRepository.On("Delete", mock.Anything, roleID).Return(model.ErrRoleNotFound)

	err := s.service.Delete(s.ctx, roleID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrRoleNotFound, err)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestDeleteRepositoryError() {
	roleID := "123e4567-e89b-12d3-a456-426614174000"

	s.roleRepository.On("Delete", mock.Anything, roleID).Return(model.ErrInternal)

	err := s.service.Delete(s.ctx, roleID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.roleRepository.AssertExpectations(s.T())
}
