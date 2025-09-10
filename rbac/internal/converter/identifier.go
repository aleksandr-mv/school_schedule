package converter

import (
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
)

// GetIdentifierToValue извлекает значение (ID или name) из GetIdentifier
// Возвращает строку, которую можно использовать для поиска по ID или имени
func GetIdentifierToValue(identifier *commonV1.GetIdentifier) (string, error) {
	if identifier == nil {
		return "", model.ErrInvalidCredentials
	}

	switch id := identifier.Identifier.(type) {
	case *commonV1.GetIdentifier_Id:
		if id.Id == "" {
			return "", model.ErrInvalidCredentials
		}
		return id.Id, nil
	case *commonV1.GetIdentifier_Name:
		if id.Name == "" {
			return "", model.ErrInvalidCredentials
		}
		return id.Name, nil
	default:
		return "", model.ErrInvalidCredentials
	}
}
