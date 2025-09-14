package role_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestListSuccess() {
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

	s.roleRepository.On("List", mock.Anything).Return(expectedRoles, nil).Once()

	roles, err := s.service.List(s.ctx)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), roles)
	assert.Len(s.T(), roles, 2)

	assert.Equal(s.T(), role1.ID, roles[0].ID)
	assert.Equal(s.T(), role1.Name, roles[0].Name)
	assert.Equal(s.T(), role1.Description, roles[0].Description)

	assert.Equal(s.T(), role2.ID, roles[1].ID)
	assert.Equal(s.T(), role2.Name, roles[1].Name)
	assert.Equal(s.T(), role2.Description, roles[1].Description)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListEmptyResult() {
	expectedRoles := []*model.Role{}

	s.roleRepository.On("List", mock.Anything).Return(expectedRoles, nil).Once()

	roles, err := s.service.List(s.ctx)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), roles)
	assert.Len(s.T(), roles, 0)

	s.roleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListInternalError() {
	s.roleRepository.On("List", mock.Anything).Return(nil, model.ErrInternal).Once()

	roles, err := s.service.List(s.ctx)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), roles)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.roleRepository.AssertExpectations(s.T())
}
