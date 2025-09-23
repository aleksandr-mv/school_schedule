package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
)

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	res, err := r.writePool.Exec(ctx, query, id)
	if err != nil {
		return r.mapDatabaseError(err, "delete")
	}

	if res.RowsAffected() == 0 {
		return model.ErrUserNotFound
	}

	return nil
}
