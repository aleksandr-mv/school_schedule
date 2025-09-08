package v1

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	authV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/auth/v1"
)

func (api *API) Logout(ctx context.Context, req *authV1.LogoutRequest) (*authV1.LogoutResponse, error) {
	sessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		logger.Warn(ctx, "❌ [API] Неверный формат UUID сессии при выходе", zap.Error(err))
		return nil, mapProtoError(ctx, model.ErrInvalidSessionData)
	}

	err = api.authService.Logout(ctx, sessionID)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка выхода из системы", zap.Error(err))
		return nil, mapProtoError(ctx, err)
	}

	return &authV1.LogoutResponse{
		Success: true,
	}, nil
}
