package model

type UserRedisView struct {
	ID                  string `redis:"user_id"`
	Login               string `redis:"user_login"`
	Email               string `redis:"user_email"`
	CreatedAtNs         int64  `redis:"user_created_at"`
	UpdatedAtNs         int64  `redis:"user_updated_at"`
	NotificationMethods string `redis:"notification_methods"`
}
