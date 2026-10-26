-- Teacher Publications: кафедры, преподаватели, категории публикаций и сами публикации.
-- Новая роль publication_manager управляет этими сущностями (по аналогии с content_manager).

-- Добавление значения в enum необратимо (PostgreSQL не умеет удалять значения enum).
-- Само значение в этой миграции нигде не используется — роль назначается позже через admins API.
ALTER TYPE admin_role ADD VALUE IF NOT EXISTS 'publication_manager';

-- ============================================================
-- DEPARTMENTS (кафедры)
-- ============================================================

CREATE TABLE departments (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    name_tm VARCHAR(255) NOT NULL,
    name_ru VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    description_tm TEXT NOT NULL DEFAULT '',
    description_ru TEXT NOT NULL DEFAULT '',
    description_en TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- TEACHERS (преподаватели)
-- ============================================================

CREATE TABLE teachers (
    id BIGSERIAL PRIMARY KEY,
    department_id BIGINT NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
    slug VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    middle_name VARCHAR(255) NOT NULL DEFAULT '',
    photo_path TEXT,
    bio TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX teachers_department_id_idx ON teachers(department_id);

-- ============================================================
-- PUBLICATION CATEGORIES (газета / журнал / конференция, динамическая таблица)
-- ============================================================

CREATE TABLE publication_categories (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    name_tm VARCHAR(255) NOT NULL,
    name_ru VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- PUBLICATIONS (статьи преподавателей)
-- ============================================================

CREATE TABLE publications (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    teacher_id BIGINT NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES publication_categories(id) ON DELETE RESTRICT,
    periodical_issue_id BIGINT REFERENCES periodical_issues(id) ON DELETE SET NULL,
    name_tm VARCHAR(512) NOT NULL,
    name_ru VARCHAR(512) NOT NULL DEFAULT '',
    name_en VARCHAR(512) NOT NULL DEFAULT '',
    article_date DATE,
    page_count INTEGER,
    file_path TEXT,
    cover_path TEXT,
    is_published BOOLEAN NOT NULL DEFAULT TRUE,
    views_count BIGINT NOT NULL DEFAULT 0,
    downloads_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX publications_teacher_id_idx ON publications(teacher_id);
CREATE INDEX publications_category_id_idx ON publications(category_id);
CREATE INDEX publications_periodical_issue_id_idx ON publications(periodical_issue_id);
CREATE INDEX publications_article_date_idx ON publications(article_date DESC);
CREATE INDEX publications_is_published_idx ON publications(is_published);

-- ============================================================
-- UPDATED_AT triggers
-- ============================================================

CREATE TRIGGER departments_updated_at BEFORE UPDATE ON departments
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER teachers_updated_at BEFORE UPDATE ON teachers
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER publication_categories_updated_at BEFORE UPDATE ON publication_categories
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER publications_updated_at BEFORE UPDATE ON publications
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ============================================================
-- SEARCH VECTOR (полнотекстовый поиск по публикациям)
-- ============================================================

ALTER TABLE publications ADD COLUMN search_vector tsvector;

CREATE OR REPLACE FUNCTION publications_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('simple', unaccent(
            coalesce(NEW.name_tm, '') || ' ' ||
            coalesce(NEW.name_ru, '') || ' ' ||
            coalesce(NEW.name_en, '')
        )), 'A');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER publications_search_vector_trigger
    BEFORE INSERT OR UPDATE OF name_tm, name_ru, name_en ON publications
    FOR EACH ROW EXECUTE FUNCTION publications_search_vector_update();

CREATE INDEX publications_search_vector_idx ON publications USING GIN(search_vector);
