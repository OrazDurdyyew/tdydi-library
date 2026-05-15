ALTER TABLE site_settings
    DROP COLUMN IF EXISTS dark_primary_color,
    DROP COLUMN IF EXISTS dark_secondary_color,
    DROP COLUMN IF EXISTS dark_accent_color,
    DROP COLUMN IF EXISTS dark_background_color,
    DROP COLUMN IF EXISTS dark_foreground_color,
    DROP COLUMN IF EXISTS dark_muted_color,
    DROP COLUMN IF EXISTS dark_border_color,
    DROP COLUMN IF EXISTS dark_ring_color;
