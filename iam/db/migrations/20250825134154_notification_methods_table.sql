-- +goose Up
-- +goose StatementBegin

-- Таблица методов уведомлений пользователей
CREATE TABLE notification_methods (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_name VARCHAR(100) NOT NULL REFERENCES providers(name) ON DELETE CASCADE,
    target VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, provider_name)
);

-- Индексы
CREATE INDEX idx_notification_methods_user_id ON notification_methods(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_methods;
-- +goose StatementEnd
