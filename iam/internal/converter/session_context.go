package converter

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/grpc/interceptor"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
)

func ExtractSessionIDFromContext(ctx context.Context) (uuid.UUID, error) {
	sessionIDStr, ok := interceptor.GetSessionIDFromContext(ctx)
	if !ok || sessionIDStr == "" {
		logger.Error(ctx, "❌ [API] Ошибка получения session ID", zap.Error(model.ErrInvalidCredentials))
		return uuid.Nil, model.ErrInvalidCredentials
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		logger.Error(ctx, "❌ [API] Невалидный session ID", zap.Error(err))
		return uuid.Nil, model.ErrInvalidCredentials
	}

	return sessionID, nil
}
