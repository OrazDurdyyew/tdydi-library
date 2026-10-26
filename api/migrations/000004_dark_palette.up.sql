-- Отдельная палитра для тёмной темы. Светлая палитра — существующие *_color поля,
-- тёмная — новые dark_*_color. Дефолты совпадают с .dark из main.css (academic).
--
-- IF NOT EXISTS: ранее эти колонки были перенесены в 000001_initial_schema.up.sql,
-- поэтому на свежей БД они уже существуют. Миграция оставлена для совместимости
-- со старыми инстансами, где initial был накатан до переноса.
ALTER TABLE site_settings
    ADD COLUMN IF NOT EXISTS dark_primary_color    VARCHAR(32) NOT NULL DEFAULT '122 162 201',
    ADD COLUMN IF NOT EXISTS dark_secondary_color  VARCHAR(32) NOT NULL DEFAULT '96 141 184',
    ADD COLUMN IF NOT EXISTS dark_accent_color     VARCHAR(32) NOT NULL DEFAULT '200 169 107',
    ADD COLUMN IF NOT EXISTS dark_background_color VARCHAR(32) NOT NULL DEFAULT '15 23 42',
    ADD COLUMN IF NOT EXISTS dark_foreground_color VARCHAR(32) NOT NULL DEFAULT '249 250 251',
    ADD COLUMN IF NOT EXISTS dark_muted_color      VARCHAR(32) NOT NULL DEFAULT '30 41 59',
    ADD COLUMN IF NOT EXISTS dark_border_color     VARCHAR(32) NOT NULL DEFAULT '31 41 55',
    ADD COLUMN IF NOT EXISTS dark_ring_color       VARCHAR(32) NOT NULL DEFAULT '122 162 201';
