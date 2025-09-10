package role

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/repository/converter"
	repoModel "github.com/aleksandr-mv/school_schedule/rbac/internal/repository/model"
)

func (r *roleRepository) Create(ctx context.Context, createRole *model.CreateRole) (*model.Role, error) {
	roleID := uuid.New()
	repoRole := converter.CreateRoleToRepo(roleID, createRole)

	query, args, err := sq.StatementBuilder.
		Insert("roles").
		Columns("id", "name", "description").
		Values(repoRole.ID, repoRole.Name, repoRole.Description).
		Suffix("RETURNING id, name, description, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to build insert query: %w", model.ErrInternal, err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // Unique constraint violation
				return nil, model.ErrRoleAlreadyExists
			case "23503": // Foreign key constraint violation
				return nil, fmt.Errorf("failed to create role: referenced entity not found")
			}
		}
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	defer rows.Close()

	createdRole, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[repoModel.Role])
	if err != nil {
		return nil, fmt.Errorf("failed to collect created role: %w", err)
	}

	return converter.RoleToDomain(&createdRole), nil
}
