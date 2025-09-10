package v1

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

// mapError преобразует внутренние ошибки в gRPC статусы
func mapError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		return status.Error(codes.InvalidArgument,
			fmt.Sprintf("Ошибка валидации: %v", validationErrors))
	}

	switch err {
	case model.ErrRoleNotFound:
		return status.Error(codes.NotFound, "Роль не найдена")
	case model.ErrRoleAlreadyExists:
		return status.Error(codes.AlreadyExists, "Роль с таким именем уже существует")
	case model.ErrInvalidRoleName:
		return status.Error(codes.InvalidArgument, "Некорректное имя роли")
	case model.ErrInvalidRoleDescription:
		return status.Error(codes.InvalidArgument, "Некорректное описание роли")
	case model.ErrInvalidCredentials:
		return status.Error(codes.InvalidArgument, "Некорректные данные")
	case model.ErrInternal:
		return status.Error(codes.Internal, "Внутренняя ошибка")
	default:
		return status.Error(codes.Internal, "Внутренняя ошибка сервера")
	}
}
