package auth_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/grpc/interceptor"
	authV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/auth/v1"
)

func (s *APISuite) TestWhoami() {
	sessionID := uuid.New()

	tests := []struct {
		name           string
		setupCtx       func() context.Context
		setupMocks     func()
		expectedErr    codes.Code
		validateResult func(*testing.T, *authV1.WhoamiResponse)
	}{
		{
			name: "Success - получение данных из Redis",
			setupCtx: func() context.Context {
				ctx := context.WithValue(s.ctx, interceptor.GetSessionIDContextKey(), sessionID.String())
				return ctx
			},
			setupMocks: func() {
				whoami := &model.WhoAMI{
					User: model.User{
						ID:    uuid.New(),
						Login: "testuser",
						Email: "test@example.com",
					},
					RolesWithPermissions: []*model.RoleWithPermissions{
						{
							Role: &model.Role{
								ID:   uuid.New(),
								Name: "admin",
							},
							Permissions: []*model.Permission{
								{ID: uuid.New(), Resource: "user", Action: "read"},
							},
						},
					},
				}
				s.whoAMIService.On("Whoami", mock.Anything, sessionID).Return(whoami, nil)
			},
			expectedErr: codes.OK,
			validateResult: func(t *testing.T, result *authV1.WhoamiResponse) {
				assert.NotNil(t, result.Info)
				assert.NotNil(t, result.Info.User)
				assert.Equal(t, "testuser", result.Info.User.Info.Login)
				assert.Equal(t, "test@example.com", result.Info.User.Info.Email)
				assert.NotEmpty(t, result.Info.RolesWithPermissions)
				assert.Equal(t, "admin", result.Info.RolesWithPermissions[0].Role.Name)
			},
		},
		{
			name: "No session ID in context",
			setupCtx: func() context.Context {
				return s.ctx
			},
			setupMocks: func() {
			},
			expectedErr: codes.Unauthenticated,
			validateResult: func(t *testing.T, result *authV1.WhoamiResponse) {
				assert.Nil(t, result)
			},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			s.whoAMIService.ExpectedCalls = nil

			ctx := tc.setupCtx()
			tc.setupMocks()

			req := &authV1.WhoamiRequest{}

			result, err := s.api.Whoami(ctx, req)

			if tc.expectedErr == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				tc.validateResult(t, result)
			} else {
				assert.Error(t, err)
				tc.validateResult(t, result)
			}

			s.whoAMIService.AssertExpectations(t)
		})
	}
}
