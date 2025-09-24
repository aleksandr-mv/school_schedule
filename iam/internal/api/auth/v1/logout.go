package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/converter"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	authV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/auth/v1"
)

func (api *API) Logout(ctx context.Context, req *authV1.LogoutRequest) (*authV1.LogoutResponse, error) {
	sessionID, err := converter.ExtractSessionIDFromContext(ctx)
	if err != nil {
		return nil, mapProtoError(ctx, err)
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
