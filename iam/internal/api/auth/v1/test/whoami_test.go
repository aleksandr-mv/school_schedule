package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	authV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/auth/v1"
)

func (s *APISuite) TestWhoami() {
	userID := uuid.New()
	sessionID := uuid.New()
	expiresAt := time.Now().Add(time.Hour)

	expectedWhoami := &model.WhoAMI{
		Session: model.Session{
			ID:        sessionID,
			ExpiresAt: expiresAt,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		User: model.User{
			ID:        userID,
			Login:     "testuser",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: nil,
			NotificationMethods: []*model.NotificationMethod{
				{
					UserID:       userID,
					ProviderName: "email",
					Target:       "test@example.com",
					CreatedAt:    time.Now(),
					UpdatedAt:    nil,
				},
			},
		},
	}

	testCases := []struct {
		name          string
		req           *authV1.WhoamiRequest
		serviceWhoami *model.WhoAMI
		serviceError  error
		expectedCode  codes.Code
		expectedError bool
	}{
		{
			name:          "Success",
			req:           &authV1.WhoamiRequest{SessionId: sessionID.String()},
			serviceWhoami: expectedWhoami,
			serviceError:  nil,
			expectedCode:  codes.OK,
		},
		{
			name:          "SessionNotFound",
			req:           &authV1.WhoamiRequest{SessionId: sessionID.String()},
			serviceWhoami: nil,
			serviceError:  model.ErrSessionNotFound,
			expectedCode:  codes.Unauthenticated,
			expectedError: true,
		},
		{
			name:          "SessionExpired",
			req:           &authV1.WhoamiRequest{SessionId: sessionID.String()},
			serviceWhoami: nil,
			serviceError:  model.ErrSessionExpired,
			expectedCode:  codes.Unauthenticated,
			expectedError: true,
		},

		{
			name:          "InternalError",
			req:           &authV1.WhoamiRequest{SessionId: sessionID.String()},
			serviceWhoami: nil,
			serviceError:  model.ErrInternal,
			expectedCode:  codes.Internal,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			s.authService.On("Whoami", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(tc.serviceWhoami, tc.serviceError).Once()

			result, err := s.api.Whoami(s.ctx, tc.req)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
				grpcErr, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedCode, grpcErr.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.Session)
				assert.NotNil(t, result.User)
				assert.Equal(t, sessionID.String(), result.Session.Id)
				assert.Equal(t, userID.String(), result.User.Id)
				assert.Equal(t, "testuser", result.User.Info.Login)
				assert.Equal(t, "test@example.com", result.User.Info.Email)
			}

			s.authService.AssertExpectations(s.T())
		})
	}
}
