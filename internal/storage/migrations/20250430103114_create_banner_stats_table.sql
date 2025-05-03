-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banner_stats (
    id SERIAL PRIMARY KEY,
    slot_id INT REFERENCES slots(id) ON DELETE CASCADE,
    banner_id INT REFERENCES banners(id) ON DELETE CASCADE,
    social_group_id INT REFERENCES social_groups(id) ON DELETE CASCADE,
    impressions INT DEFAULT 0,
    clicks INT DEFAULT 0,
    UNIQUE(slot_id, banner_id, social_group_id)
);

COMMENT ON TABLE banner_stats IS 'Статистика показов и кликов баннеров';

COMMENT ON COLUMN banner_stats.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN banner_stats.slot_id IS 'ID слота';
COMMENT ON COLUMN banner_stats.banner_id IS 'ID баннера';
COMMENT ON COLUMN banner_stats.social_group_id IS 'ID соц-дем группы';
COMMENT ON COLUMN banner_stats.impressions IS 'Число показов';
COMMENT ON COLUMN banner_stats.clicks IS 'Число кликов';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banner_stats;
-- +goose StatementEnd
