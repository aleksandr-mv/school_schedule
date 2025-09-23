package v1

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
)

func mapProtoError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, model.ErrInvalidCredentials):
		return status.Errorf(codes.Unauthenticated, "invalid login or password")
	case errors.Is(err, model.ErrSessionExpired):
		return status.Errorf(codes.Unauthenticated, "session expired")
	case errors.Is(err, model.ErrSessionNotFound):
		return status.Errorf(codes.Unauthenticated, "session not found")

	case errors.Is(err, model.ErrUserNotFound):
		return status.Errorf(codes.NotFound, "user not found")

	case errors.Is(err, model.ErrNotificationNotFound):
		return status.Errorf(codes.NotFound, "notification method not found")

	case errors.Is(err, model.ErrInvalidSessionData):
		return status.Errorf(codes.InvalidArgument, "invalid session data")

	case errors.Is(err, model.ErrFailedToCreateSession),
		errors.Is(err, model.ErrFailedToDeleteSession),
		errors.Is(err, model.ErrFailedToStoreInCache),
		errors.Is(err, model.ErrFailedToReadFromCache),
		errors.Is(err, model.ErrFailedToGetUser),
		errors.Is(err, model.ErrFailedToListNotifications),
		errors.Is(err, model.ErrInternal):
		return status.Errorf(codes.Internal, "internal server error")
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return status.Errorf(codes.InvalidArgument, "validation error: %s", validationErrs.Error())
	}

	logger.Error(ctx, "❌ [API] Неожиданная ошибка", zap.Error(err))
	return status.Errorf(codes.Internal, "internal server error")
}
