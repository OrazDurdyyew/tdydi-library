-- Initial schema for электронная библиотека
-- Расширения (pg_trgm, unaccent) подключаются init.sql на старте PostgreSQL.

-- ============================================================
-- ADMINS + PERMISSIONS
-- ============================================================

CREATE TYPE admin_role AS ENUM ('superadmin', 'content_manager');

CREATE TABLE admins (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email VARCHAR(255),
    role admin_role NOT NULL DEFAULT 'content_manager',
    all_sections BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX admins_role_idx ON admins(role);

-- ============================================================
-- SITE SETTINGS (singleton: только одна строка с id=1)
-- ============================================================

CREATE TABLE site_settings (
    id SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    name_tm VARCHAR(255) NOT NULL DEFAULT '',
    name_ru VARCHAR(255) NOT NULL DEFAULT '',
    name_en VARCHAR(255) NOT NULL DEFAULT '',
    logo_light_path TEXT,
    logo_dark_path TEXT,
    favicon_path TEXT,
    primary_color VARCHAR(32) NOT NULL DEFAULT '30 58 95',
    secondary_color VARCHAR(32) NOT NULL DEFAULT '44 95 141',
    accent_color VARCHAR(32) NOT NULL DEFAULT '200 169 107',
    background_color VARCHAR(32) NOT NULL DEFAULT '248 250 252',
    foreground_color VARCHAR(32) NOT NULL DEFAULT '17 24 39',
    muted_color VARCHAR(32) NOT NULL DEFAULT '241 245 249',
    border_color VARCHAR(32) NOT NULL DEFAULT '229 231 235',
    ring_color VARCHAR(32) NOT NULL DEFAULT '44 95 141',
    dark_primary_color VARCHAR(32) NOT NULL DEFAULT '122 162 201',
    dark_secondary_color VARCHAR(32) NOT NULL DEFAULT '96 141 184',
    dark_accent_color VARCHAR(32) NOT NULL DEFAULT '200 169 107',
    dark_background_color VARCHAR(32) NOT NULL DEFAULT '15 23 42',
    dark_foreground_color VARCHAR(32) NOT NULL DEFAULT '249 250 251',
    dark_muted_color VARCHAR(32) NOT NULL DEFAULT '30 41 59',
    dark_border_color VARCHAR(32) NOT NULL DEFAULT '31 41 55',
    dark_ring_color VARCHAR(32) NOT NULL DEFAULT '122 162 201',
    font_family VARCHAR(255) NOT NULL DEFAULT 'Inter',
    contact_address TEXT,
    contact_phone VARCHAR(64),
    contact_email VARCHAR(255),
    social_links JSONB NOT NULL DEFAULT '{}'::jsonb,
    footer_tm TEXT NOT NULL DEFAULT '',
    footer_ru TEXT NOT NULL DEFAULT '',
    footer_en TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO site_settings (id) VALUES (1);

-- ============================================================
-- SECTIONS (2-уровневая иерархия: раздел → подраздел)
-- ============================================================

CREATE TABLE sections (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    parent_id BIGINT REFERENCES sections(id) ON DELETE RESTRICT,
    sort_order INTEGER NOT NULL DEFAULT 0,
    image_path TEXT,
    name_tm VARCHAR(255) NOT NULL,
    name_ru VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    description_tm TEXT NOT NULL DEFAULT '',
    description_ru TEXT NOT NULL DEFAULT '',
    description_en TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX sections_parent_id_idx ON sections(parent_id);
CREATE INDEX sections_sort_order_idx ON sections(sort_order);

-- Триггер: запрещает >2 уровня вложенности (подраздел не может иметь подразделов)
CREATE OR REPLACE FUNCTION sections_check_max_depth()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.parent_id IS NOT NULL THEN
        IF EXISTS (
            SELECT 1 FROM sections WHERE id = NEW.parent_id AND parent_id IS NOT NULL
        ) THEN
            RAISE EXCEPTION 'Sections support only 2 levels of nesting (root -> child)';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER sections_max_depth_trigger
    BEFORE INSERT OR UPDATE OF parent_id ON sections
    FOR EACH ROW EXECUTE FUNCTION sections_check_max_depth();

-- ============================================================
-- ADMIN SECTION PERMISSIONS (M:N: какие разделы доступны админу)
-- ============================================================

CREATE TABLE admin_section_permissions (
    admin_id BIGINT NOT NULL REFERENCES admins(id) ON DELETE CASCADE,
    section_id BIGINT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    PRIMARY KEY (admin_id, section_id)
);

CREATE INDEX admin_section_permissions_section_id_idx ON admin_section_permissions(section_id);

-- ============================================================
-- BOOKS
-- ============================================================

CREATE TABLE books (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    isbn VARCHAR(32),
    year SMALLINT,
    pages INTEGER,
    language VARCHAR(8) NOT NULL DEFAULT 'tm',
    title VARCHAR(512) NOT NULL,
    annotation TEXT NOT NULL DEFAULT '',
    authors TEXT NOT NULL DEFAULT '',
    publisher VARCHAR(255),
    file_path TEXT,
    file_format VARCHAR(16),
    cover_path TEXT,
    allow_download BOOLEAN NOT NULL DEFAULT TRUE,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    views_count BIGINT NOT NULL DEFAULT 0,
    downloads_count BIGINT NOT NULL DEFAULT 0,
    created_by_admin_id BIGINT REFERENCES admins(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX books_is_published_idx ON books(is_published);
CREATE INDEX books_year_idx ON books(year);
CREATE INDEX books_language_idx ON books(language);
CREATE INDEX books_created_at_idx ON books(created_at DESC);

-- ============================================================
-- BOOK_SECTIONS (M:N: книга может быть в нескольких разделах/подразделах)
-- ============================================================

CREATE TABLE book_sections (
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    section_id BIGINT NOT NULL REFERENCES sections(id) ON DELETE RESTRICT,
    PRIMARY KEY (book_id, section_id)
);

CREATE INDEX book_sections_section_id_idx ON book_sections(section_id);

-- ============================================================
-- TAGS + BOOK_TAGS
-- ============================================================

CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    name_tm VARCHAR(128) NOT NULL,
    name_ru VARCHAR(128) NOT NULL,
    name_en VARCHAR(128) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE book_tags (
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, tag_id)
);

CREATE INDEX book_tags_tag_id_idx ON book_tags(tag_id);

-- ============================================================
-- COLLECTIONS (виртуальные подборки)
-- ============================================================

CREATE TABLE collections (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    cover_path TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0,
    show_on_home BOOLEAN NOT NULL DEFAULT FALSE,
    name_tm VARCHAR(255) NOT NULL,
    name_ru VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    description_tm TEXT NOT NULL DEFAULT '',
    description_ru TEXT NOT NULL DEFAULT '',
    description_en TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX collections_show_on_home_idx ON collections(show_on_home);

CREATE TABLE collection_books (
    collection_id BIGINT NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (collection_id, book_id)
);

CREATE INDEX collection_books_book_id_idx ON collection_books(book_id);

-- ============================================================
-- PERIODICALS (газеты + журналы в одной сущности)
-- ============================================================

CREATE TYPE periodical_type AS ENUM ('newspaper', 'magazine');

CREATE TABLE periodicals (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    type periodical_type NOT NULL,
    cover_path TEXT,
    periodicity VARCHAR(64),
    name_tm VARCHAR(255) NOT NULL,
    name_ru VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    description_tm TEXT NOT NULL DEFAULT '',
    description_ru TEXT NOT NULL DEFAULT '',
    description_en TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX periodicals_type_idx ON periodicals(type);

CREATE TABLE periodical_issues (
    id BIGSERIAL PRIMARY KEY,
    periodical_id BIGINT NOT NULL REFERENCES periodicals(id) ON DELETE CASCADE,
    issue_number VARCHAR(64),
    issue_date DATE NOT NULL,
    file_path TEXT,
    cover_path TEXT,
    annotation_tm TEXT NOT NULL DEFAULT '',
    annotation_ru TEXT NOT NULL DEFAULT '',
    annotation_en TEXT NOT NULL DEFAULT '',
    views_count BIGINT NOT NULL DEFAULT 0,
    downloads_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX periodical_issues_periodical_id_idx ON periodical_issues(periodical_id);
CREATE INDEX periodical_issues_date_idx ON periodical_issues(issue_date DESC);

-- ============================================================
-- POSTS (новости/блог)
-- ============================================================

CREATE TYPE post_status AS ENUM ('draft', 'published');

CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    cover_path TEXT,
    author VARCHAR(255),
    published_at TIMESTAMPTZ,
    status post_status NOT NULL DEFAULT 'draft',
    title_tm VARCHAR(512) NOT NULL,
    title_ru VARCHAR(512) NOT NULL,
    title_en VARCHAR(512) NOT NULL,
    content_tm TEXT NOT NULL DEFAULT '',
    content_ru TEXT NOT NULL DEFAULT '',
    content_en TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX posts_status_idx ON posts(status);
CREATE INDEX posts_published_at_idx ON posts(published_at DESC);

-- ============================================================
-- PAGES (статичные страницы CMS)
-- ============================================================

CREATE TABLE pages (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    in_menu BOOLEAN NOT NULL DEFAULT FALSE,
    menu_order INTEGER NOT NULL DEFAULT 0,
    title_tm VARCHAR(255) NOT NULL,
    title_ru VARCHAR(255) NOT NULL,
    title_en VARCHAR(255) NOT NULL,
    content_tm TEXT NOT NULL DEFAULT '',
    content_ru TEXT NOT NULL DEFAULT '',
    content_en TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX pages_in_menu_idx ON pages(in_menu);

-- ============================================================
-- FAQ
-- ============================================================

CREATE TABLE faq_categories (
    id BIGSERIAL PRIMARY KEY,
    sort_order INTEGER NOT NULL DEFAULT 0,
    name_tm VARCHAR(255) NOT NULL,
    name_ru VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE faq_items (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES faq_categories(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    question_tm TEXT NOT NULL,
    question_ru TEXT NOT NULL,
    question_en TEXT NOT NULL,
    answer_tm TEXT NOT NULL,
    answer_ru TEXT NOT NULL,
    answer_en TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX faq_items_category_id_idx ON faq_items(category_id);

-- ============================================================
-- HOME PAGE CONFIGURATION
-- ============================================================

CREATE TABLE home_slides (
    id BIGSERIAL PRIMARY KEY,
    image_path TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    link_url TEXT,
    title_tm VARCHAR(255) NOT NULL DEFAULT '',
    title_ru VARCHAR(255) NOT NULL DEFAULT '',
    title_en VARCHAR(255) NOT NULL DEFAULT '',
    subtitle_tm TEXT NOT NULL DEFAULT '',
    subtitle_ru TEXT NOT NULL DEFAULT '',
    subtitle_en TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE home_sections (
    id BIGSERIAL PRIMARY KEY,
    section_id BIGINT NOT NULL UNIQUE REFERENCES sections(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE home_collections (
    id BIGSERIAL PRIMARY KEY,
    collection_id BIGINT NOT NULL UNIQUE REFERENCES collections(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- ============================================================
-- VIEW LOGS (для статистики просмотров)
-- ============================================================

CREATE TYPE view_entity_type AS ENUM ('book', 'periodical_issue', 'section', 'post');

CREATE TABLE view_logs (
    id BIGSERIAL PRIMARY KEY,
    entity_type view_entity_type NOT NULL,
    entity_id BIGINT NOT NULL,
    ip_hash VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX view_logs_entity_idx ON view_logs(entity_type, entity_id);
CREATE INDEX view_logs_created_at_idx ON view_logs(created_at DESC);

-- ============================================================
-- UPDATED_AT auto-trigger
-- ============================================================

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER admins_updated_at BEFORE UPDATE ON admins
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER site_settings_updated_at BEFORE UPDATE ON site_settings
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER sections_updated_at BEFORE UPDATE ON sections
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER books_updated_at BEFORE UPDATE ON books
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER collections_updated_at BEFORE UPDATE ON collections
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER periodicals_updated_at BEFORE UPDATE ON periodicals
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER periodical_issues_updated_at BEFORE UPDATE ON periodical_issues
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER posts_updated_at BEFORE UPDATE ON posts
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER pages_updated_at BEFORE UPDATE ON pages
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER faq_categories_updated_at BEFORE UPDATE ON faq_categories
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER faq_items_updated_at BEFORE UPDATE ON faq_items
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
