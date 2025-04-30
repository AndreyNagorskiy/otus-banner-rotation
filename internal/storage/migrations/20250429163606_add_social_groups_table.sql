-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS social_groups (
   id SERIAL PRIMARY KEY,
   description TEXT NOT NULL
);

COMMENT ON TABLE social_groups IS 'Группа пользователей сайта со схожими интересами, например "девушки 20-25" или "дедушки 80+"';

COMMENT ON COLUMN social_groups.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN social_groups.description IS 'Описание';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS social_groups;
-- +goose StatementEnd
