-- Rollback search vectors

-- PAGES
DROP INDEX IF EXISTS pages_search_vector_idx;
DROP TRIGGER IF EXISTS pages_search_vector_trigger ON pages;
DROP FUNCTION IF EXISTS pages_search_vector_update();
ALTER TABLE pages DROP COLUMN IF EXISTS search_vector;

-- PERIODICALS
DROP INDEX IF EXISTS periodicals_search_vector_idx;
DROP TRIGGER IF EXISTS periodicals_search_vector_trigger ON periodicals;
DROP FUNCTION IF EXISTS periodicals_search_vector_update();
ALTER TABLE periodicals DROP COLUMN IF EXISTS search_vector;

-- FAQ ITEMS
DROP INDEX IF EXISTS faq_items_search_vector_idx;
DROP TRIGGER IF EXISTS faq_items_search_vector_trigger ON faq_items;
DROP FUNCTION IF EXISTS faq_items_search_vector_update();
ALTER TABLE faq_items DROP COLUMN IF EXISTS search_vector;

-- POSTS
DROP INDEX IF EXISTS posts_search_vector_idx;
DROP TRIGGER IF EXISTS posts_search_vector_trigger ON posts;
DROP FUNCTION IF EXISTS posts_search_vector_update();
ALTER TABLE posts DROP COLUMN IF EXISTS search_vector;

-- BOOKS
DROP INDEX IF EXISTS books_title_trgm_idx;
DROP INDEX IF EXISTS books_search_vector_idx;
DROP TRIGGER IF EXISTS books_search_vector_trigger ON books;
DROP FUNCTION IF EXISTS books_search_vector_update();
ALTER TABLE books DROP COLUMN IF EXISTS search_vector;
