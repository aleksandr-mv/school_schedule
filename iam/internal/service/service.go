package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
)

type UserService interface {
	Register(ctx context.Context, login, email, password string) (*model.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type AuthService interface {
	Login(ctx context.Context, credentials *model.LoginCredentials) (uuid.UUID, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
}

type UserProducerService interface {
	ProduceUserCreated(ctx context.Context, event model.UserCreated) error
}

type WhoAMIService interface {
	Whoami(ctx context.Context, sessionID uuid.UUID) (*model.WhoAMI, error)
}
