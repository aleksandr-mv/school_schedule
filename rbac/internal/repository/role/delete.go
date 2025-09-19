package role

import (
	"context"
	"fmt"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (r *roleRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE roles SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.writePool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: delete role failed: %w", model.ErrInternal, err)
	}

	if result.RowsAffected() == 0 {
		return model.ErrRoleNotFound
	}

	return nil
}
