-- Search vectors + GIN indexes для полнотекстового поиска.
-- Используем конфиг 'simple' (lowercase + tokenize, без стемминга), так как контент на 3 языках.
-- unaccent убирает диакритику; pg_trgm даёт нечёткий поиск (опечатки).

-- ============================================================
-- BOOKS
-- ============================================================

ALTER TABLE books ADD COLUMN search_vector tsvector;

CREATE OR REPLACE FUNCTION books_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('simple', unaccent(coalesce(NEW.title, ''))), 'A') ||
        setweight(to_tsvector('simple', unaccent(coalesce(NEW.authors, ''))), 'B') ||
        setweight(to_tsvector('simple', unaccent(coalesce(NEW.annotation, ''))), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER books_search_vector_trigger
    BEFORE INSERT OR UPDATE OF title, authors, annotation ON books
    FOR EACH ROW EXECUTE FUNCTION books_search_vector_update();

CREATE INDEX books_search_vector_idx ON books USING GIN(search_vector);
CREATE INDEX books_title_trgm_idx ON books USING GIN(title gin_trgm_ops);

-- ============================================================
-- POSTS (новости)
-- ============================================================

ALTER TABLE posts ADD COLUMN search_vector tsvector;

CREATE OR REPLACE FUNCTION posts_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('simple', unaccent(
            coalesce(NEW.title_tm, '') || ' ' ||
            coalesce(NEW.title_ru, '') || ' ' ||
            coalesce(NEW.title_en, '')
        )), 'A') ||
        setweight(to_tsvector('simple', unaccent(
            coalesce(NEW.content_tm, '') || ' ' ||
            coalesce(NEW.content_ru, '') || ' ' ||
            coalesce(NEW.content_en, '')
        )), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER posts_search_vector_trigger
    BEFORE INSERT OR UPDATE OF title_tm, title_ru, title_en, content_tm, content_ru, content_en ON posts
    FOR EACH ROW EXECUTE FUNCTION posts_search_vector_update();

CREATE INDEX posts_search_vector_idx ON posts USING GIN(search_vector);

-- ============================================================
-- FAQ ITEMS
-- ============================================================

ALTER TABLE faq_items ADD COLUMN search_vector tsvector;

CREATE OR REPLACE FUNCTION faq_items_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('simple', unaccent(
            coalesce(NEW.question_tm, '') || ' ' ||
            coalesce(NEW.question_ru, '') || ' ' ||
            coalesce(NEW.question_en, '')
        )), 'A') ||
        setweight(to_tsvector('simple', unaccent(
            coalesce(NEW.answer_tm, '') || ' ' ||
            coalesce(NEW.answer_ru, '') || ' ' ||
            coalesce(NEW.answer_en, '')
        )), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER faq_items_search_vector_trigger
    BEFORE INSERT OR UPDATE OF question_tm, question_ru, question_en, answer_tm, answer_ru, answer_en ON faq_items
    FOR EACH ROW EXECUTE FUNCTION faq_items_search_vector_update();

CREATE INDEX faq_items_search_vector_idx ON faq_items USING GIN(search_vector);

-- ============================================================
-- PERIODICALS
-- ============================================================

ALTER TABLE periodicals ADD COLUMN search_vector tsvector;

CREATE OR REPLACE FUNCTION periodicals_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('simple', unaccent(
            coalesce(NEW.name_tm, '') || ' ' ||
            coalesce(NEW.name_ru, '') || ' ' ||
            coalesce(NEW.name_en, '')
        )), 'A') ||
        setweight(to_tsvector('simple', unaccent(
            coalesce(NEW.description_tm, '') || ' ' ||
            coalesce(NEW.description_ru, '') || ' ' ||
            coalesce(NEW.description_en, '')
        )), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER periodicals_search_vector_trigger
    BEFORE INSERT OR UPDATE OF name_tm, name_ru, name_en, description_tm, description_ru, description_en ON periodicals
    FOR EACH ROW EXECUTE FUNCTION periodicals_search_vector_update();

CREATE INDEX periodicals_search_vector_idx ON periodicals USING GIN(search_vector);

-- ============================================================
-- PAGES
-- ============================================================

ALTER TABLE pages ADD COLUMN search_vector tsvector;

CREATE OR REPLACE FUNCTION pages_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('simple', unaccent(
            coalesce(NEW.title_tm, '') || ' ' ||
            coalesce(NEW.title_ru, '') || ' ' ||
            coalesce(NEW.title_en, '')
        )), 'A') ||
        setweight(to_tsvector('simple', unaccent(
            coalesce(NEW.content_tm, '') || ' ' ||
            coalesce(NEW.content_ru, '') || ' ' ||
            coalesce(NEW.content_en, '')
        )), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER pages_search_vector_trigger
    BEFORE INSERT OR UPDATE OF title_tm, title_ru, title_en, content_tm, content_ru, content_en ON pages
    FOR EACH ROW EXECUTE FUNCTION pages_search_vector_update();

CREATE INDEX pages_search_vector_idx ON pages USING GIN(search_vector);
