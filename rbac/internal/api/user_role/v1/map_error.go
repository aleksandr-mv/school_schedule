package v1

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, model.ErrUserRoleNotFound):
		return status.Error(codes.NotFound, "Связь пользователь-роль не найдена")
	case errors.Is(err, model.ErrRoleNotFound):
		return status.Error(codes.NotFound, "Роль не найдена")
	case errors.Is(err, model.ErrRoleAlreadyAssigned):
		return status.Error(codes.AlreadyExists, "Роль уже назначена пользователю")
	case errors.Is(err, model.ErrRoleNotAssigned):
		return status.Error(codes.FailedPrecondition, "Роль не назначена пользователю")
	case errors.Is(err, model.ErrInternal):
		return status.Error(codes.Internal, "Внутренняя ошибка")
	default:
		return status.Error(codes.Internal, "Внутренняя ошибка сервера")
	}
}
