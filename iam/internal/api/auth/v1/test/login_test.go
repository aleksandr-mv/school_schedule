package auth_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	authV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/auth/v1"
)

func (s *APISuite) TestLogin() {
	expectedSessionID := uuid.New()

	testCases := []struct {
		name             string
		req              *authV1.LoginRequest
		serviceSessionID uuid.UUID
		serviceError     error
		expectedCode     codes.Code
		expectedError    bool
	}{
		{
			name: "Success",
			req: &authV1.LoginRequest{
				Login:    "testuser",
				Password: "password123",
			},
			serviceSessionID: expectedSessionID,
			serviceError:     nil,
			expectedCode:     codes.OK,
		},
		{
			name: "InvalidCredentials",
			req: &authV1.LoginRequest{
				Login:    "testuser",
				Password: "wrongpassword",
			},
			serviceSessionID: uuid.Nil,
			serviceError:     model.ErrInvalidCredentials,
			expectedCode:     codes.Unauthenticated,
			expectedError:    true,
		},
		{
			name: "UserNotFound",
			req: &authV1.LoginRequest{
				Login:    "nonexistent",
				Password: "password123",
			},
			serviceSessionID: uuid.Nil,
			serviceError:     model.ErrInvalidCredentials,
			expectedCode:     codes.Unauthenticated,
			expectedError:    true,
		},
		{
			name: "InternalError",
			req: &authV1.LoginRequest{
				Login:    "testuser",
				Password: "password123",
			},
			serviceSessionID: uuid.Nil,
			serviceError:     model.ErrInternal,
			expectedCode:     codes.Internal,
			expectedError:    true,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			s.authService.On("Login", mock.Anything, mock.AnythingOfType("*model.LoginCredentials")).Return(tc.serviceSessionID, tc.serviceError).Once()

			result, err := s.api.Login(s.ctx, tc.req)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
				grpcErr, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedCode, grpcErr.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, expectedSessionID.String(), result.SessionId)
			}

			s.authService.AssertExpectations(s.T())
		})
	}
}
