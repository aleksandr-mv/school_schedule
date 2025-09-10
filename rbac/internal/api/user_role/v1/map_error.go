package v1

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func mapError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		return status.Error(codes.InvalidArgument,
			fmt.Sprintf("Ошибка валидации: %v", validationErrors))
	}

	switch err {
	case model.ErrUserRoleNotFound:
		return status.Error(codes.NotFound, "Связь пользователь-роль не найдена")
	case model.ErrRoleNotFound:
		return status.Error(codes.NotFound, "Роль не найдена")
	case model.ErrRoleAlreadyAssigned:
		return status.Error(codes.AlreadyExists, "Роль уже назначена пользователю")
	case model.ErrRoleNotAssigned:
		return status.Error(codes.FailedPrecondition, "Роль не назначена пользователю")
	case model.ErrInvalidCredentials:
		return status.Error(codes.InvalidArgument, "Некорректные данные")
	case model.ErrInternal:
		return status.Error(codes.Internal, "Внутренняя ошибка")
	default:
		return status.Error(codes.Internal, "Внутренняя ошибка сервера")
	}
}
