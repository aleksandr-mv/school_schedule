package notification

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
)

func (r *notificationRepository) mapDatabaseError(err error, operation string) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return model.ErrNotificationAlreadyExists
		case "23503":
			return model.ErrNotificationUserConstraintViolation
		case "23502":
			return model.ErrInvalidNotificationData
		default:
			return fmt.Errorf("database constraint violation (code: %s): %w", pgErr.Code, err)
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return model.ErrNotificationNotFound
	}

	switch operation {
	case "create":
		return fmt.Errorf("%w: %w", model.ErrFailedToCreateNotification, err)
	case "delete":
		return fmt.Errorf("%w: %w", model.ErrFailedToDeleteNotification, err)
	case "get", "select":
		return fmt.Errorf("%w: %w", model.ErrFailedToGetNotification, err)
	case "list":
		return fmt.Errorf("%w: %w", model.ErrFailedToListNotifications, err)
	default:
		return fmt.Errorf("notification repository operation failed: %w", err)
	}
}
