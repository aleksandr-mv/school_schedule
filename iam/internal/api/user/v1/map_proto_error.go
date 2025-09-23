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
	case errors.Is(err, model.ErrUserNotFound):
		return status.Errorf(codes.NotFound, "user not found")
	case errors.Is(err, model.ErrUserAlreadyExists):
		return status.Errorf(codes.AlreadyExists, "user already exists")

	case errors.Is(err, model.ErrUserConstraintViolation):
		return status.Errorf(codes.InvalidArgument, "user constraint violation")

	case errors.Is(err, model.ErrFailedToGetNotification):
		return status.Errorf(codes.Internal, "failed to get notification method")
	case errors.Is(err, model.ErrFailedToListNotifications):
		return status.Errorf(codes.Internal, "failed to list notification methods")

	case errors.Is(err, model.ErrInvalidSessionData):
		return status.Errorf(codes.InvalidArgument, "invalid session data")

	case errors.Is(err, model.ErrFailedToCreateUser),
		errors.Is(err, model.ErrFailedToUpdateUser),
		errors.Is(err, model.ErrFailedToDeleteUser),
		errors.Is(err, model.ErrFailedToGetUser),
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
