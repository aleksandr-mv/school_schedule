package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/converter"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	authV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/auth/v1"
)

func (api *API) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	sessionID, err := api.authService.Login(ctx, converter.LoginFromProto(req))
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка входа в систему", zap.Error(err))
		return nil, mapProtoError(ctx, err)
	}

	logger.Info(ctx, "✅ [API] Пользователь успешно вошел в систему")
	return &authV1.LoginResponse{
		SessionId: sessionID.String(),
	}, nil
}
