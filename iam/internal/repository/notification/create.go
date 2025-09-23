package notification

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/converter"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/model"
)

func (r *notificationRepository) Create(ctx context.Context, notificationMethod model.NotificationMethod) (*model.NotificationMethod, error) {
	repoNotificationMethod := converter.NotificationMethodToRepo(&notificationMethod)

	query, args, err := sq.StatementBuilder.
		Insert("notification_methods").
		Columns("user_id", "provider_name", "target").
		Values(repoNotificationMethod.UserID, repoNotificationMethod.ProviderName, repoNotificationMethod.Target).
		Suffix("RETURNING user_id, provider_name, target, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to build insert query: %w", model.ErrInternal, err)
	}

	rows, err := r.writePool.Query(ctx, query, args...)
	if err != nil {
		return nil, r.mapDatabaseError(err, "create")
	}
	defer rows.Close()

	createdNotification, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[repoModel.NotificationMethod])
	if err != nil {
		return nil, r.mapDatabaseError(err, "create")
	}

	return converter.NotificationMethodFromRepo(&createdNotification), nil
}
