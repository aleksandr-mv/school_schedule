package converter

import (
	"time"

	"github.com/google/uuid"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	repoModel "github.com/aleksandr-mv/school_schedule/iam/internal/repository/model"
)

func UserToRedis(user *model.User) (repoModel.UserRedisView, error) {
	notificationMethodsJSON, err := marshalNotificationMethods(user.NotificationMethods)
	if err != nil {
		return repoModel.UserRedisView{}, err
	}

	var updatedAtNs int64
	if user.UpdatedAt != nil {
		updatedAtNs = user.UpdatedAt.UnixNano()
	}

	return repoModel.UserRedisView{
		ID:                  user.ID.String(),
		Login:               user.Login,
		Email:               user.Email,
		CreatedAtNs:         user.CreatedAt.UnixNano(),
		UpdatedAtNs:         updatedAtNs,
		NotificationMethods: notificationMethodsJSON,
	}, nil
}

func UserFromRedis(redisView repoModel.UserRedisView) (*model.User, error) {
	id, err := uuid.Parse(redisView.ID)
	if err != nil {
		return nil, err
	}

	notificationMethods, err := unmarshalNotificationMethods(redisView.NotificationMethods)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:                  id,
		Login:               redisView.Login,
		Email:               redisView.Email,
		CreatedAt:           time.Unix(0, redisView.CreatedAtNs),
		NotificationMethods: notificationMethods,
	}

	if redisView.UpdatedAtNs > 0 {
		updatedAt := time.Unix(0, redisView.UpdatedAtNs)
		user.UpdatedAt = &updatedAt
	}

	return user, nil
}
