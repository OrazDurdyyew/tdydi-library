-- Rollback teacher publications.
-- Примечание: значение 'publication_manager' в enum admin_role не удаляется —
-- PostgreSQL не поддерживает удаление значений enum (необратимо, no-op).

DROP INDEX IF EXISTS publications_search_vector_idx;
DROP TRIGGER IF EXISTS publications_search_vector_trigger ON publications;
DROP FUNCTION IF EXISTS publications_search_vector_update();

DROP TRIGGER IF EXISTS publications_updated_at ON publications;
DROP TRIGGER IF EXISTS publication_categories_updated_at ON publication_categories;
DROP TRIGGER IF EXISTS teachers_updated_at ON teachers;
DROP TRIGGER IF EXISTS departments_updated_at ON departments;

DROP TABLE IF EXISTS publications;
DROP TABLE IF EXISTS publication_categories;
DROP TABLE IF EXISTS teachers;
DROP TABLE IF EXISTS departments;
