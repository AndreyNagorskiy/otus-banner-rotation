-- +goose Up
-- +goose StatementBegin
INSERT INTO slots (description) VALUES
    ('Главный баннер на главной странице'),
    ('Боковая панель в каталоге'),
    ('Баннер в футере');

INSERT INTO banners (description) VALUES
    ('Скидка 50% на всё'),
    ('Новая коллекция весна-лето'),
    ('Подпишитесь на рассылку');

INSERT INTO social_groups (description) VALUES
    ('Женщины 20-30 лет'),
    ('Мужчины 30-40 лет'),
    ('Пенсионеры 60+');

-- Привязка баннеров к слотам
INSERT INTO slot_banners (slot_id, banner_id) VALUES
    (1, 1),
    (1, 2),
    (2, 2),
    (2, 3),
    (3, 1),
    (3, 3);

-- Заполнение статистики
INSERT INTO banner_stats (slot_id, banner_id, social_group_id, impressions, clicks) VALUES
    (1, 1, 1, 10, 3),
    (1, 2, 1, 8, 1),
    (2, 2, 2, 5, 2),
    (2, 3, 2, 6, 0),
    (3, 1, 3, 7, 4),
    (3, 3, 3, 3, 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE banners, slots, social_groups, slot_banners, banner_stats CASCADE;
-- +goose StatementEnd
