package converter

import (
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/model"
)

func ToRepoNotificationMethod(method *model.NotificationMethod) *repoModel.NotificationMethod {
	return &repoModel.NotificationMethod{
		UserID:       method.UserID,
		ProviderName: method.ProviderName,
		Target:       method.Target,
		CreatedAt:    method.CreatedAt,
		UpdatedAt:    method.UpdatedAt,
	}
}

func ToDomainNotificationMethod(method *repoModel.NotificationMethod) *model.NotificationMethod {
	return &model.NotificationMethod{
		UserID:       method.UserID,
		ProviderName: method.ProviderName,
		Target:       method.Target,
		CreatedAt:    method.CreatedAt,
		UpdatedAt:    method.UpdatedAt,
	}
}

func ToDomainNotificationMethods(methods []repoModel.NotificationMethod) []*model.NotificationMethod {
	result := make([]*model.NotificationMethod, len(methods))
	for i, method := range methods {
		result[i] = ToDomainNotificationMethod(&method)
	}
	return result
}
