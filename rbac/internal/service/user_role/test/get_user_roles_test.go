package user_role_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestGetUserRolesSuccess() {
	userID := "user123"
	roleID1 := uuid.New().String()
	roleID2 := uuid.New().String()

	role1 := &model.EnrichedRole{
		Role: model.Role{
			ID:        uuid.New(),
			Name:      "admin",
			CreatedAt: time.Now(),
		},
		Permissions: []*model.Permission{},
	}
	role2 := &model.EnrichedRole{
		Role: model.Role{
			ID:        uuid.New(),
			Name:      "user",
			CreatedAt: time.Now(),
		},
		Permissions: []*model.Permission{},
	}

	s.userRoleRepository.On("GetUserRoles", mock.Anything, userID).Return([]string{roleID1, roleID2}, nil)
	s.roleService.On("Get", mock.Anything, roleID1).Return(role1, nil)
	s.roleService.On("Get", mock.Anything, roleID2).Return(role2, nil)

	result, err := s.service.GetUserRoles(s.ctx, userID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 2)

	assert.Equal(s.T(), role1.Role.ID, result[0].Role.ID)
	assert.Equal(s.T(), role1.Role.Name, result[0].Role.Name)
	assert.Equal(s.T(), role2.Role.ID, result[1].Role.ID)
	assert.Equal(s.T(), role2.Role.Name, result[1].Role.Name)

	s.userRoleRepository.AssertExpectations(s.T())
	s.roleService.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetUserRolesEmptyResult() {
	userID := "user123"

	s.userRoleRepository.On("GetUserRoles", mock.Anything, userID).Return([]string{}, nil)

	result, err := s.service.GetUserRoles(s.ctx, userID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 0)

	s.userRoleRepository.AssertExpectations(s.T())
}
