package user_role_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	userRoleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user_role/v1"
)

func (s *APISuite) TestAssignSuccess() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	req := &userRoleV1.AssignRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: &assignedBy,
	}

	s.userRoleService.On("Assign", mock.Anything, userID, roleID, &assignedBy).Return(nil).Once()

	resp, err := s.api.Assign(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignSuccessWithoutAssignedBy() {
	userID := "user123"
	roleID := "role456"

	req := &userRoleV1.AssignRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: nil,
	}

	s.userRoleService.On("Assign", mock.Anything, userID, roleID, (*string)(nil)).Return(nil).Once()

	resp, err := s.api.Assign(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignRoleAlreadyAssigned() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	req := &userRoleV1.AssignRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: &assignedBy,
	}

	s.userRoleService.On("Assign", mock.Anything, userID, roleID, &assignedBy).Return(model.ErrRoleAlreadyAssigned).Once()

	resp, err := s.api.Assign(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.AlreadyExists, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignUserNotFound() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	req := &userRoleV1.AssignRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: &assignedBy,
	}

	s.userRoleService.On("Assign", mock.Anything, userID, roleID, &assignedBy).Return(model.ErrRoleNotFound).Once()

	resp, err := s.api.Assign(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignRoleNotFound() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	req := &userRoleV1.AssignRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: &assignedBy,
	}

	s.userRoleService.On("Assign", mock.Anything, userID, roleID, &assignedBy).Return(model.ErrRoleNotFound).Once()

	resp, err := s.api.Assign(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignInternalError() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	req := &userRoleV1.AssignRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: &assignedBy,
	}

	s.userRoleService.On("Assign", mock.Anything, userID, roleID, &assignedBy).Return(model.ErrInternal).Once()

	resp, err := s.api.Assign(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignValidation_InvalidUserID() {
	req := &userRoleV1.AssignRequest{
		UserId: "invalid-uuid",
		RoleId: uuid.New().String(),
	}

	err := req.Validate()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "value must be a valid UUID")
}

func (s *APISuite) TestAssignValidation_InvalidRoleID() {
	req := &userRoleV1.AssignRequest{
		UserId: uuid.New().String(),
		RoleId: "invalid-uuid",
	}

	err := req.Validate()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "value must be a valid UUID")
}

func (s *APISuite) TestAssignValidation_InvalidAssignedBy() {
	assignedBy := "invalid-uuid"

	req := &userRoleV1.AssignRequest{
		UserId:     uuid.New().String(),
		RoleId:     uuid.New().String(),
		AssignedBy: &assignedBy,
	}

	err := req.Validate()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "value must be a valid UUID")
}

func (s *APISuite) TestAssignValidation_ValidRequest() {
	userID := uuid.New().String()
	roleID := uuid.New().String()
	assignedBy := uuid.New().String()

	req := &userRoleV1.AssignRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: &assignedBy,
	}

	err := req.Validate()
	assert.NoError(s.T(), err)
}

func (s *APISuite) TestAssignValidation_ValidRequestWithoutAssignedBy() {
	userID := uuid.New().String()
	roleID := uuid.New().String()

	req := &userRoleV1.AssignRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: nil,
	}

	err := req.Validate()
	assert.NoError(s.T(), err)
}
