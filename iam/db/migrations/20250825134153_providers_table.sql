-- +goose Up
-- +goose StatementBegin

-- Справочник провайдеров уведомлений
CREATE TABLE providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Вставка начальных данных для провайдеров
INSERT INTO providers (name, description) VALUES
    ('telegram', 'Уведомления через Telegram бота'),
    ('email', 'Уведомления по электронной почте'),
    ('push', 'Push уведомления в браузере'),
    ('sms', 'SMS уведомления');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS providers;
-- +goose StatementEnd