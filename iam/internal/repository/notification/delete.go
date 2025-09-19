package notification

import (
	"context"

	"github.com/google/uuid"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
)

func (r *notificationRepository) Delete(ctx context.Context, userID uuid.UUID, providerName string) error {
	query := `DELETE FROM notification_methods WHERE user_id = $1 AND provider_name = $2`

	res, err := r.writePool.Exec(ctx, query, userID, providerName)
	if err != nil {
		return r.mapDatabaseError(err, "delete")
	}

	if res.RowsAffected() == 0 {
		return model.ErrNotificationNotFound
	}

	return nil
}
