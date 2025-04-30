-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banners (
     id SERIAL PRIMARY KEY,
     description TEXT NOT NULL
);

COMMENT ON TABLE banners IS 'Рекламный/информационный элемент, который показывается в слоте';

COMMENT ON COLUMN banners.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN banners.description IS 'Описание';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banners;
-- +goose StatementEnd
