package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/converter"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	authV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/auth/v1"
)

func (api *API) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	sessionID, err := converter.ExtractSessionIDFromContext(ctx)
	if err != nil {
		return nil, mapProtoError(ctx, err)
	}

	whoami, err := api.whoAMIService.Whoami(ctx, sessionID)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения данных пользователя из Redis", zap.Error(err))
		return nil, mapProtoError(ctx, err)
	}

	return &authV1.WhoamiResponse{
		Info: converter.WhoamiToProto(whoami),
	}, nil
}
