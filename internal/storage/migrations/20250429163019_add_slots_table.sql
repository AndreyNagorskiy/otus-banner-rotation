-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS slots (
   id SERIAL PRIMARY KEY,
   description TEXT NOT NULL
);

COMMENT ON TABLE slots IS 'Место на сайте, на котором показывается баннер';

COMMENT ON COLUMN slots.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN slots.description IS 'Описание';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS slots;
-- +goose StatementEnd
