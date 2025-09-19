package v1_test

import (
	"time"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
)

func (s *APISuite) TestCheckSuccess() {
	sessionID := uuid.New()
	userID := uuid.New()

	// Подготавливаем данные WhoAMI
	whoami := &model.WhoAMI{
		Session: model.Session{
			ID:        sessionID,
			ExpiresAt: time.Now().Add(time.Hour),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		User: model.User{
			ID:        userID,
			Login:     "testuser",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
		},
		RolesWithPermissions: []*model.RoleWithPermissions{
			{
				Role: &model.Role{
					ID:   uuid.New(),
					Name: "admin",
				},
				Permissions: []*model.Permission{
					{
						ID:       uuid.New(),
						Resource: "users",
						Action:   "read",
					},
					{
						ID:       uuid.New(),
						Resource: "users",
						Action:   "write",
					},
				},
			},
		},
	}

	// Подготавливаем запрос с session UUID в заголовке
	req := &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Request: &authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Headers: map[string]string{
						"session-uuid": sessionID.String(),
					},
				},
			},
		},
	}

	// Настраиваем мок
	s.whoAMIService.On("Whoami", mock.Anything, sessionID).Return(whoami, nil)

	// Выполняем запрос
	result, err := s.api.Check(s.ctx, req)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), int32(0), result.Status.Code) // Успешный статус

	// Проверяем что это OkResponse
	okResponse, ok := result.HttpResponse.(*authv3.CheckResponse_OkResponse)
	assert.True(s.T(), ok)
	assert.NotNil(s.T(), okResponse.OkResponse)

	// Проверяем заголовки
	headers := okResponse.OkResponse.Headers
	assert.Len(s.T(), headers, 4) // UUID, Login, Roles, Permissions

	// Проверяем конкретные заголовки
	headerMap := make(map[string]string)
	for _, header := range headers {
		headerMap[header.Header.Key] = header.Header.Value
	}

	assert.Equal(s.T(), userID.String(), headerMap["X-User-UUID"])
	assert.Equal(s.T(), "testuser", headerMap["X-User-Login"])
	assert.Equal(s.T(), "admin", headerMap["X-User-Roles"])
	assert.Contains(s.T(), headerMap["X-User-Permissions"], "users:read")
	assert.Contains(s.T(), headerMap["X-User-Permissions"], "users:write")

	s.whoAMIService.AssertExpectations(s.T())
}

func (s *APISuite) TestCheckMissingSession() {
	// Запрос без session UUID
	req := &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Request: &authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Headers: map[string]string{},
				},
			},
		},
	}

	// Выполняем запрос
	result, err := s.api.Check(s.ctx, req)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), int32(16), result.Status.Code) // Unauthenticated

	// Проверяем что это DeniedResponse
	deniedResponse, ok := result.HttpResponse.(*authv3.CheckResponse_DeniedResponse)
	assert.True(s.T(), ok)
	assert.NotNil(s.T(), deniedResponse.DeniedResponse)
	assert.Contains(s.T(), deniedResponse.DeniedResponse.Body, "Missing or invalid session")
}

func (s *APISuite) TestCheckInvalidSession() {
	sessionID := uuid.New()

	// Подготавливаем запрос
	req := &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Request: &authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Headers: map[string]string{
						"session-uuid": sessionID.String(),
					},
				},
			},
		},
	}

	// Настраиваем мок на ошибку
	s.whoAMIService.On("Whoami", mock.Anything, sessionID).Return(nil, model.ErrSessionNotFound)

	// Выполняем запрос
	result, err := s.api.Check(s.ctx, req)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), int32(16), result.Status.Code) // Unauthenticated

	// Проверяем что это DeniedResponse
	deniedResponse, ok := result.HttpResponse.(*authv3.CheckResponse_DeniedResponse)
	assert.True(s.T(), ok)
	assert.NotNil(s.T(), deniedResponse.DeniedResponse)
	assert.Contains(s.T(), deniedResponse.DeniedResponse.Body, "Invalid session")

	s.whoAMIService.AssertExpectations(s.T())
}
