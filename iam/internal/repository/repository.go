package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
	Get(ctx context.Context, value string) (*model.User, error)
	Update(ctx context.Context, user model.User) (*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type NotificationRepository interface {
	Create(ctx context.Context, notificationMethod model.NotificationMethod) (*model.NotificationMethod, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*model.NotificationMethod, error)
	GetByUserAndProvider(ctx context.Context, userID uuid.UUID, providerName string) (*model.NotificationMethod, error)
	Delete(ctx context.Context, userID uuid.UUID, providerName string) error
}

type SessionRepository interface {
	Create(ctx context.Context, user model.User, expiresAt time.Time) (uuid.UUID, error)
	Get(ctx context.Context, sessionID uuid.UUID) (*model.WhoAMI, error)
	Delete(ctx context.Context, sessionID uuid.UUID) error
}
