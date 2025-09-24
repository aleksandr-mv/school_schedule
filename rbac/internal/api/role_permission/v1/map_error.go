package v1

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, model.ErrPermissionNotFound):
		return status.Error(codes.NotFound, "Право доступа не найдено")
	case errors.Is(err, model.ErrRoleNotFound):
		return status.Error(codes.NotFound, "Роль не найдена")
	case errors.Is(err, model.ErrRolePermissionNotFound):
		return status.Error(codes.NotFound, "Связь роль-право не найдена")
	case errors.Is(err, model.ErrPermissionAlreadyAssigned):
		return status.Error(codes.AlreadyExists, "Право уже назначено роли")
	case errors.Is(err, model.ErrPermissionNotAssigned):
		return status.Error(codes.FailedPrecondition, "Право не назначено роли")
	case errors.Is(err, model.ErrInternal):
		return status.Error(codes.Internal, "Внутренняя ошибка")
	default:
		return status.Error(codes.Internal, "Внутренняя ошибка сервера")
	}
}
