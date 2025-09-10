package user_role_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestGetRoleUsersSuccess() {
	roleID := "role123"
	limit := int32(10)
	cursor := "cursor123"
	nextCursor := "nextCursor456"

	expectedUsers := []string{"user1", "user2", "user3"}
	expectedTotal := int32(3)

	s.userRoleRepository.On("GetRoleUsers", mock.Anything, roleID, limit, &cursor).Return(expectedUsers, expectedTotal, &nextCursor, nil)

	users, total, nextCursorResult, err := s.service.GetRoleUsers(s.ctx, roleID, limit, &cursor)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 3)
	assert.Equal(s.T(), expectedUsers, users)
	assert.Equal(s.T(), expectedTotal, total)
	assert.NotNil(s.T(), nextCursorResult)
	assert.Equal(s.T(), nextCursor, *nextCursorResult)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetRoleUsersWithoutCursor() {
	roleID := "role123"
	limit := int32(10)

	expectedUsers := []string{"user1", "user2"}
	expectedTotal := int32(2)

	s.userRoleRepository.On("GetRoleUsers", mock.Anything, roleID, limit, (*string)(nil)).Return(expectedUsers, expectedTotal, nil, nil)

	users, total, nextCursorResult, err := s.service.GetRoleUsers(s.ctx, roleID, limit, nil)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)
	assert.Equal(s.T(), expectedUsers, users)
	assert.Equal(s.T(), expectedTotal, total)
	assert.Nil(s.T(), nextCursorResult)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetRoleUsersEmptyResult() {
	roleID := "role123"
	limit := int32(10)

	s.userRoleRepository.On("GetRoleUsers", mock.Anything, roleID, limit, (*string)(nil)).Return([]string{}, int32(0), nil, nil)

	users, total, nextCursorResult, err := s.service.GetRoleUsers(s.ctx, roleID, limit, nil)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 0)
	assert.Equal(s.T(), int32(0), total)
	assert.Nil(s.T(), nextCursorResult)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetRoleUsersRepositoryError() {
	roleID := "role123"
	limit := int32(10)

	s.userRoleRepository.On("GetRoleUsers", mock.Anything, roleID, limit, (*string)(nil)).Return(nil, int32(0), nil, model.ErrInternal)

	users, total, nextCursorResult, err := s.service.GetRoleUsers(s.ctx, roleID, limit, nil)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)
	assert.Nil(s.T(), users)
	assert.Equal(s.T(), int32(0), total)
	assert.Nil(s.T(), nextCursorResult)

	s.userRoleRepository.AssertExpectations(s.T())
}
