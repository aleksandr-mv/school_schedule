package user_role

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *userRoleRepository) GetRoleUsers(ctx context.Context, roleID string, limit int32, cursor *string) ([]string, int32, *string, error) {
	// Получаем общее количество пользователей с этой ролью
	countQuery := `SELECT COUNT(*) FROM user_roles WHERE role_id = $1`
	var totalCount int32
	err := r.pool.QueryRow(ctx, countQuery, roleID).Scan(&totalCount)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to count role users: %w", err)
	}

	// Строим запрос для получения пользователей с пагинацией
	query := `SELECT user_id FROM user_roles WHERE role_id = $1`
	args := []interface{}{roleID}
	argIndex := 2

	if cursor != nil && *cursor != "" {
		query += ` AND user_id > $` + fmt.Sprintf("%d", argIndex)
		args = append(args, *cursor)
		argIndex++
	}

	query += ` ORDER BY user_id LIMIT $` + fmt.Sprintf("%d", argIndex)
	args = append(args, limit+1) // +1 для проверки наличия следующей страницы

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to get role users: %w", err)
	}
	defer rows.Close()

	userIDs, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to collect role users: %w", err)
	}

	// Проверяем, есть ли следующая страница
	var nextCursor *string
	if len(userIDs) > int(limit) {
		// Убираем лишний элемент и устанавливаем курсор
		userIDs = userIDs[:limit]
		nextCursor = &userIDs[len(userIDs)-1]
	}

	return userIDs, totalCount, nextCursor, nil
}
