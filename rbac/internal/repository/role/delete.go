package role

import (
	"context"
	"fmt"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (r *roleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM roles WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return model.ErrRoleNotFound
	}

	return nil
}
