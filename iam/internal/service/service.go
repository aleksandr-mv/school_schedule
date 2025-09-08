package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
)

type UserServiceInterface interface {
	Register(ctx context.Context, createUser *model.CreateUser) (*model.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type AuthServiceInterface interface {
	Login(ctx context.Context, credentials *model.LoginCredentials) (uuid.UUID, error)
	Whoami(ctx context.Context, sessionID uuid.UUID) (*model.WhoAMI, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
}
