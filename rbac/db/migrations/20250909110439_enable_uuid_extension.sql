-- +goose Up
-- +goose StatementBegin
-- Включаем расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Отключаем расширение UUID (осторожно - может сломать другие таблицы)
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd