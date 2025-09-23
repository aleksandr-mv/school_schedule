package user

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
)

func (r *userRepository) mapDatabaseError(err error, operation string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return model.ErrUserNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("%w: %s", model.ErrUserAlreadyExists, pgErr.Detail)
		case "23503":
			return model.ErrUserConstraintViolation
		case "23502":
			return model.ErrInvalidUserData
		case "22P02":
			return fmt.Errorf("%w: invalid id format", model.ErrInvalidUserData)
		default:
			return fmt.Errorf("database constraint violation (%s): %w", pgErr.Code, err)
		}
	}

	switch operation {
	case "create":
		return fmt.Errorf("%w: %w", model.ErrFailedToCreateUser, err)
	case "update":
		return fmt.Errorf("%w: %w", model.ErrFailedToUpdateUser, err)
	case "delete":
		return fmt.Errorf("%w: %w", model.ErrFailedToDeleteUser, err)
	case "get":
		return fmt.Errorf("%w: %w", model.ErrFailedToGetUser, err)
	case "exist":
		return fmt.Errorf("%w: %w", model.ErrFailedToGetUser, err)
	default:
		return fmt.Errorf("user repository operation failed: %w", err)
	}
}
