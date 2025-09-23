package user_role_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestGetUsersSuccess() {
	roleID := "role123"
	limit := int32(10)
	cursor := "cursor123"
	nextCursor := "nextCursor456"

	expectedUsers := []string{"user1", "user2", "user3"}

	s.userRoleRepository.On("GetRoleUsers", mock.Anything, roleID, limit, cursor).Return(expectedUsers, &nextCursor, nil)

	users, nextCursorResult, err := s.service.GetRoleUsers(s.ctx, roleID, limit, cursor)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 3)
	assert.Equal(s.T(), expectedUsers, users)
	assert.NotNil(s.T(), nextCursorResult)
	assert.Equal(s.T(), nextCursor, *nextCursorResult)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetUsersWithoutCursor() {
	roleID := "role123"
	limit := int32(10)

	expectedUsers := []string{"user1", "user2"}

	s.userRoleRepository.On("GetRoleUsers", mock.Anything, roleID, limit, "").Return(expectedUsers, nil, nil)

	users, nextCursorResult, err := s.service.GetRoleUsers(s.ctx, roleID, limit, "")

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)
	assert.Equal(s.T(), expectedUsers, users)
	assert.Nil(s.T(), nextCursorResult)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetUsersEmptyResult() {
	roleID := "role123"
	limit := int32(10)

	s.userRoleRepository.On("GetRoleUsers", mock.Anything, roleID, limit, "").Return([]string{}, nil, nil)

	users, nextCursorResult, err := s.service.GetRoleUsers(s.ctx, roleID, limit, "")

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 0)
	assert.Nil(s.T(), nextCursorResult)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetUsersRepositoryError() {
	roleID := "role123"
	limit := int32(10)

	s.userRoleRepository.On("GetRoleUsers", mock.Anything, roleID, limit, "").Return(nil, nil, model.ErrInternal)

	users, nextCursorResult, err := s.service.GetRoleUsers(s.ctx, roleID, limit, "")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)
	assert.Nil(s.T(), users)
	assert.Nil(s.T(), nextCursorResult)

	s.userRoleRepository.AssertExpectations(s.T())
}
