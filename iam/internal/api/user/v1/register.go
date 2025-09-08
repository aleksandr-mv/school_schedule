package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/iam/internal/converter"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	userV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user/v1"
)

func (api *API) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	user, err := api.userService.Register(ctx, converter.UserFromProto(req))
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка регистрации пользователя", zap.Error(err))
		return nil, mapProtoError(ctx, err)
	}

	logger.Info(ctx, "✅ [API] Пользователь успешно зарегистрирован")
	return &userV1.RegisterResponse{
		UserId: user.ID.String(),
	}, nil
}
