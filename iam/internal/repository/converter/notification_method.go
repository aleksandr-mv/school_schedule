package converter

import (
	"encoding/json"
	"fmt"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	repoModel "github.com/aleksandr-mv/school_schedule/iam/internal/repository/model"
)

func NotificationMethodToRepo(notificationMethod *model.NotificationMethod) *repoModel.NotificationMethod {
	repoNotificationMethod := &repoModel.NotificationMethod{
		UserID:       notificationMethod.UserID,
		ProviderName: notificationMethod.ProviderName,
		Target:       notificationMethod.Target,
		CreatedAt:    notificationMethod.CreatedAt,
	}

	if notificationMethod.UpdatedAt != nil {
		repoNotificationMethod.UpdatedAt = notificationMethod.UpdatedAt
	}

	return repoNotificationMethod
}

func NotificationMethodFromRepo(notificationMethod *repoModel.NotificationMethod) *model.NotificationMethod {
	domainNotificationMethod := &model.NotificationMethod{
		UserID:       notificationMethod.UserID,
		ProviderName: notificationMethod.ProviderName,
		Target:       notificationMethod.Target,
		CreatedAt:    notificationMethod.CreatedAt,
	}

	if notificationMethod.UpdatedAt != nil {
		domainNotificationMethod.UpdatedAt = notificationMethod.UpdatedAt
	}

	return domainNotificationMethod
}

func NotificationMethodListFromRepo(notificationMethods []repoModel.NotificationMethod) ([]*model.NotificationMethod, error) {
	result := make([]*model.NotificationMethod, len(notificationMethods))
	for i, nm := range notificationMethods {
		converted := NotificationMethodFromRepo(&nm)
		if converted == nil {
			return nil, fmt.Errorf("failed to convert notification method at index %d", i)
		}
		result[i] = converted
	}
	return result, nil
}

func marshalNotificationMethods(methods []*model.NotificationMethod) (string, error) {
	if methods == nil {
		return "", nil
	}

	bytes, err := json.Marshal(methods)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func unmarshalNotificationMethods(jsonStr string) ([]*model.NotificationMethod, error) {
	if jsonStr == "" {
		return nil, nil
	}

	var tempMethods []model.NotificationMethod
	if err := json.Unmarshal([]byte(jsonStr), &tempMethods); err != nil {
		return nil, err
	}

	notificationMethods := make([]*model.NotificationMethod, len(tempMethods))
	for i := range tempMethods {
		notificationMethods[i] = &tempMethods[i]
	}

	return notificationMethods, nil
}
