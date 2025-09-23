package auth_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	authV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/auth/v1"
)

func (s *APISuite) TestLogout() {
	sessionID := uuid.New()

	testCases := []struct {
		name          string
		req           *authV1.LogoutRequest
		serviceError  error
		expectedCode  codes.Code
		expectedError bool
	}{
		{
			name:          "Success",
			req:           &authV1.LogoutRequest{SessionId: sessionID.String()},
			serviceError:  nil,
			expectedCode:  codes.OK,
			expectedError: false,
		},
		{
			name:          "SessionNotFound",
			req:           &authV1.LogoutRequest{SessionId: sessionID.String()},
			serviceError:  model.ErrSessionNotFound,
			expectedCode:  codes.Unauthenticated,
			expectedError: true,
		},

		{
			name:          "InternalError",
			req:           &authV1.LogoutRequest{SessionId: sessionID.String()},
			serviceError:  model.ErrInternal,
			expectedCode:  codes.Internal,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			s.authService.On("Logout", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(tc.serviceError).Once()

			result, err := s.api.Logout(s.ctx, tc.req)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
				grpcErr, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedCode, grpcErr.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			s.authService.AssertExpectations(s.T())
		})
	}
}
