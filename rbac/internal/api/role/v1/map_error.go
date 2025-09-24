package v1

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

// mapError преобразует внутренние ошибки в gRPC статусы
func mapError(err error) error {
	switch {
	case errors.Is(err, model.ErrRoleNotFound):
		return status.Error(codes.NotFound, "Роль не найдена")
	case errors.Is(err, model.ErrRoleAlreadyExists):
		return status.Error(codes.AlreadyExists, "Роль с таким именем уже существует")
	case errors.Is(err, model.ErrFailedToCreateRole):
		return status.Error(codes.Internal, "Не удалось создать роль")
	case errors.Is(err, model.ErrInternal):
		return status.Error(codes.Internal, "Внутренняя ошибка")
	default:
		return status.Error(codes.Internal, "Внутренняя ошибка сервера")
	}
}
