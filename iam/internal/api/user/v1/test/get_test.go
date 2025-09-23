package user_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	userV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/user/v1"
)

func (s *APISuite) TestGetUser() {
	userID := uuid.New()

	expectedUser := &model.User{
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
	}

	testCases := []struct {
		name          string
		req           *userV1.GetUserRequest
		serviceUser   *model.User
		serviceError  error
		expectedCode  codes.Code
		expectedError bool
	}{
		{
			name:         "Success",
			req:          &userV1.GetUserRequest{UserId: userID.String()},
			serviceUser:  expectedUser,
			serviceError: nil,
			expectedCode: codes.OK,
		},
		{
			name:          "UserNotFound",
			req:           &userV1.GetUserRequest{UserId: userID.String()},
			serviceUser:   nil,
			serviceError:  model.ErrUserNotFound,
			expectedCode:  codes.NotFound,
			expectedError: true,
		},

		{
			name:          "InternalError",
			req:           &userV1.GetUserRequest{UserId: userID.String()},
			serviceUser:   nil,
			serviceError:  model.ErrInternal,
			expectedCode:  codes.Internal,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			s.userService.On("GetUser", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(tc.serviceUser, tc.serviceError).Once()

			result, err := s.api.GetUser(s.ctx, tc.req)
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
				grpcErr, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedCode, grpcErr.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.User)
				assert.Equal(t, userID.String(), result.User.Id)
				assert.Equal(t, "testuser", result.User.Info.Login)
				assert.Equal(t, "test@example.com", result.User.Info.Email)
				assert.Len(t, result.User.Info.NotificationMethods, 1)
			}

			s.userService.AssertExpectations(s.T())
		})
	}
}
