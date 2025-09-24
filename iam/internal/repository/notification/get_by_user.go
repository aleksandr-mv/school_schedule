package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/converter"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/model"
)

func (r *notificationRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*model.NotificationMethod, error) {
	query := `SELECT user_id, provider_name, target, created_at, updated_at 
			  FROM notification_methods 
			  WHERE user_id = $1 
			  ORDER BY created_at ASC`

	rows, err := r.readPool.Query(ctx, query, userID)
	if err != nil {
		return nil, r.mapDatabaseError(err, "list")
	}
	defer rows.Close()

	notificationMethods, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[repoModel.NotificationMethod])
	if err != nil {
		return nil, r.mapDatabaseError(err, "list")
	}

	return converter.ToDomainNotificationMethods(notificationMethods), nil
}
