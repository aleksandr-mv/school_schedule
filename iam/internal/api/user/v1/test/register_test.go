package user_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
	userV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user/v1"
)

func (s *APISuite) TestRegister() {
	userID := uuid.New()

	expectedUser := &model.User{
		ID:        userID,
		Login:     "newuser",
		Email:     "new@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	testCases := []struct {
		name          string
		req           *userV1.RegisterRequest
		serviceUser   *model.User
		serviceError  error
		expectedCode  codes.Code
		expectedError bool
	}{
		{
			name: "Success",
			req: &userV1.RegisterRequest{
				Info: &commonV1.UserInfo{
					Login: "newuser",
					Email: "new@example.com",
				},
				Password: "password123",
			},
			serviceUser:  expectedUser,
			serviceError: nil,
			expectedCode: codes.OK,
		},
		{
			name: "UserAlreadyExists",
			req: &userV1.RegisterRequest{
				Info: &commonV1.UserInfo{
					Login: "existinguser",
					Email: "existing@example.com",
				},
				Password: "password123",
			},
			serviceUser:   nil,
			serviceError:  model.ErrUserAlreadyExists,
			expectedCode:  codes.AlreadyExists,
			expectedError: true,
		},
		{
			name: "InternalError",
			req: &userV1.RegisterRequest{
				Info: &commonV1.UserInfo{
					Login: "testuser",
					Email: "test@example.com",
				},
				Password: "password123",
			},
			serviceUser:   nil,
			serviceError:  model.ErrInternal,
			expectedCode:  codes.Internal,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			s.userService.On("Register", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(tc.serviceUser, tc.serviceError).Once()

			result, err := s.api.Register(s.ctx, tc.req)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
				grpcErr, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedCode, grpcErr.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, userID.String(), result.UserId)
			}

			s.userService.AssertExpectations(s.T())
		})
	}
}
