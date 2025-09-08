package v1

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/iam/internal/converter"
	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	authV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/auth/v1"
)

func (api *API) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	sessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		logger.Warn(ctx, "❌ [API] Неверный формат UUID сессии", zap.Error(err))
		return nil, mapProtoError(ctx, model.ErrInvalidSessionData)
	}

	iam, err := api.authService.Whoami(ctx, sessionID)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения информации о сессии", zap.Error(err))
		return nil, mapProtoError(ctx, err)
	}

	return converter.WhoAMIToProto(iam), nil
}
