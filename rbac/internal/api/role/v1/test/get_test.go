package role_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
)

func (s *APISuite) TestGetRoleSuccessByID() {
	roleID := uuid.New()
	expectedRole := &model.Role{
		ID:          roleID,
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}

	req := &roleV1.GetRoleRequest{
		Value: &commonV1.GetIdentifier{
			Identifier: &commonV1.GetIdentifier_Id{
				Id: roleID.String(),
			},
		},
	}

	s.roleService.On("GetRole", mock.Anything, mock.Anything).Return(expectedRole, nil).Once()

	resp, err := s.api.GetRole(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotNil(s.T(), resp.Role)
	assert.Equal(s.T(), expectedRole.ID.String(), resp.Role.Id)
	assert.Equal(s.T(), expectedRole.Name, resp.Role.Name)
	assert.Equal(s.T(), expectedRole.Description, resp.Role.Description)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetRoleSuccessByName() {
	roleID := uuid.New()
	expectedRole := &model.Role{
		ID:          roleID,
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}

	req := &roleV1.GetRoleRequest{
		Value: &commonV1.GetIdentifier{
			Identifier: &commonV1.GetIdentifier_Name{
				Name: "admin",
			},
		},
	}

	s.roleService.On("GetRole", mock.Anything, mock.Anything).Return(expectedRole, nil).Once()

	resp, err := s.api.GetRole(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotNil(s.T(), resp.Role)
	assert.Equal(s.T(), expectedRole.ID.String(), resp.Role.Id)
	assert.Equal(s.T(), expectedRole.Name, resp.Role.Name)
	assert.Equal(s.T(), expectedRole.Description, resp.Role.Description)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetRoleInvalidIdentifier() {
	req := &roleV1.GetRoleRequest{
		Value: &commonV1.GetIdentifier{
			Identifier: nil, // пустой идентификатор
		},
	}

	// Сервис не вызывается при ошибке конвертера
	resp, err := s.api.GetRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.InvalidArgument, grpcErr.Code())
}

func (s *APISuite) TestGetRoleNotFound() {
	roleID := uuid.New()

	req := &roleV1.GetRoleRequest{
		Value: &commonV1.GetIdentifier{
			Identifier: &commonV1.GetIdentifier_Id{
				Id: roleID.String(),
			},
		},
	}

	s.roleService.On("GetRole", mock.Anything, mock.Anything).Return(nil, model.ErrRoleNotFound).Once()

	resp, err := s.api.GetRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetRoleInternalError() {
	roleID := uuid.New()

	req := &roleV1.GetRoleRequest{
		Value: &commonV1.GetIdentifier{
			Identifier: &commonV1.GetIdentifier_Id{
				Id: roleID.String(),
			},
		},
	}

	s.roleService.On("GetRole", mock.Anything, mock.Anything).Return(nil, model.ErrInternal).Once()

	resp, err := s.api.GetRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}
