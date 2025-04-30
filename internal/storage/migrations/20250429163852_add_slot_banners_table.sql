-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS slot_banners (
    id SERIAL PRIMARY KEY,
    slot_id INT NOT NULL REFERENCES slots(id),
    banner_id INT NOT NULL REFERENCES banners(id),
    UNIQUE(slot_id, banner_id)
);

COMMENT ON TABLE slot_banners IS 'Связь между слотами и баннерами, участвующими в ротации';

COMMENT ON COLUMN slot_banners.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN slot_banners.slot_id IS 'ID слота, в котором участвует баннер';
COMMENT ON COLUMN slot_banners.banner_id IS 'ID баннера, участвующего в слоте';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS slot_banners;
-- +goose StatementEnd
