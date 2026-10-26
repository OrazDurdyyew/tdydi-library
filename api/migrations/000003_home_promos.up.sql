-- Промо-блоки главной страницы: картинка + заголовок/подзаголовок (3 языка) + ссылка.
CREATE TABLE home_promos (
    id BIGSERIAL PRIMARY KEY,
    image_path TEXT NOT NULL DEFAULT '',
    link_url TEXT NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    title_tm VARCHAR(255) NOT NULL DEFAULT '',
    title_ru VARCHAR(255) NOT NULL DEFAULT '',
    title_en VARCHAR(255) NOT NULL DEFAULT '',
    subtitle_tm TEXT NOT NULL DEFAULT '',
    subtitle_ru TEXT NOT NULL DEFAULT '',
    subtitle_en TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX home_promos_active_idx ON home_promos (is_active, sort_order);

-- Академическая палитра по умолчанию (navy + gold). Обновляем только установки,
-- где палитра ещё дефолтная — кастомизированные ВУЗами цвета не трогаем.
UPDATE site_settings SET
    primary_color = '30 58 95',
    secondary_color = '44 95 141',
    accent_color = '200 169 107',
    background_color = '248 250 252',
    foreground_color = '17 24 39',
    muted_color = '241 245 249',
    border_color = '229 231 235',
    ring_color = '44 95 141'
WHERE id = 1 AND primary_color = '22 171 101';
