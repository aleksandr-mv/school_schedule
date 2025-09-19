package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	"github.com/aleksandr-mv/school_schedule/iam/internal/repository/converter"
	repoModel "github.com/aleksandr-mv/school_schedule/iam/internal/repository/model"
)

func (r *notificationRepository) GetByUserAndProvider(ctx context.Context, userID uuid.UUID, providerName string) (*model.NotificationMethod, error) {
	query := `SELECT user_id, provider_name, target, created_at, updated_at 
			  FROM notification_methods 
			  WHERE user_id = $1 AND provider_name = $2`

	rows, err := r.readPool.Query(ctx, query, userID, providerName)
	if err != nil {
		return nil, r.mapDatabaseError(err, "get")
	}
	defer rows.Close()

	notificationMethod, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[repoModel.NotificationMethod])
	if err != nil {
		return nil, r.mapDatabaseError(err, "get")
	}

	return converter.NotificationMethodFromRepo(&notificationMethod), nil
}
