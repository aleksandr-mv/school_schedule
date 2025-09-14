package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func mapError(err error) error {
	switch err {
	case model.ErrPermissionNotFound:
		return status.Error(codes.NotFound, "Право доступа не найдено")
	case model.ErrRoleNotFound:
		return status.Error(codes.NotFound, "Роль не найдена")
	case model.ErrRolePermissionNotFound:
		return status.Error(codes.NotFound, "Связь роль-право не найдена")
	case model.ErrPermissionAlreadyAssigned:
		return status.Error(codes.AlreadyExists, "Право уже назначено роли")
	case model.ErrPermissionNotAssigned:
		return status.Error(codes.FailedPrecondition, "Право не назначено роли")
	case model.ErrInternal:
		return status.Error(codes.Internal, "Внутренняя ошибка")
	default:
		return status.Error(codes.Internal, "Внутренняя ошибка сервера")
	}
}
