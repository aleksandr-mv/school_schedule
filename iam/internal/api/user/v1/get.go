package v1

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/converter"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	userV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/user/v1"
)

func (api *API) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Warn(ctx, "❌ [API] Неверный формат UUID пользователя", zap.Error(err))
		return nil, mapProtoError(ctx, model.ErrInvalidSessionData)
	}

	user, err := api.userService.GetUser(ctx, userID)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения пользователя", zap.Error(err))
		return nil, mapProtoError(ctx, err)
	}

	return &userV1.GetUserResponse{
		User: converter.UserToProto(user),
	}, nil
}
