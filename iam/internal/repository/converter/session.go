package converter

import (
	"time"

	"github.com/google/uuid"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	repoModel "github.com/aleksandr-mv/school_schedule/iam/internal/repository/model"
)

func SessionToRedis(session model.Session) repoModel.SessionRedisView {
	return repoModel.SessionRedisView{
		ID:          session.ID.String(),
		ExpiresAtNs: session.ExpiresAt.UnixNano(),
		CreatedAtNs: session.CreatedAt.UnixNano(),
		UpdatedAtNs: session.UpdatedAt.UnixNano(),
	}
}

func SessionFromRedis(redisView repoModel.SessionRedisView) (model.Session, error) {
	id, err := uuid.Parse(redisView.ID)
	if err != nil {
		return model.Session{}, err
	}

	return model.Session{
		ID:        id,
		ExpiresAt: time.Unix(0, redisView.ExpiresAtNs),
		CreatedAt: time.Unix(0, redisView.CreatedAtNs),
		UpdatedAt: time.Unix(0, redisView.UpdatedAtNs),
	}, nil
}

func CreateSessionCacheView(session model.Session, user *model.User) (repoModel.SessionCacheView, error) {
	sessionView := SessionToRedis(session)

	userView, err := UserToRedis(user)
	if err != nil {
		return repoModel.SessionCacheView{}, err
	}

	return repoModel.SessionCacheView{
		SessionRedisView: sessionView,
		UserRedisView:    userView,
	}, nil
}
