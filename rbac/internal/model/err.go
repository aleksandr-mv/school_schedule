package model

import "errors"

// Ошибки доменной модели
var (
	ErrRoleNotFound              = errors.New("роль не найдена")
	ErrRoleAlreadyExists         = errors.New("роль с таким именем уже существует")
	ErrPermissionNotFound        = errors.New("право доступа не найдено")
	ErrUserRoleNotFound          = errors.New("связь пользователь-роль не найдена")
	ErrRolePermissionNotFound    = errors.New("связь роль-право не найдена")
	ErrPermissionAlreadyAssigned = errors.New("право уже назначено роли")
	ErrPermissionNotAssigned     = errors.New("право не назначено роли")
	ErrRoleAlreadyAssigned       = errors.New("роль уже назначена пользователю")
	ErrRoleNotAssigned           = errors.New("роль не назначена пользователю")
	ErrFailedToCreateRole        = errors.New("не удалось создать роль")
	ErrInternal                  = errors.New("внутренняя ошибка")
)
