DROP INDEX IF EXISTS books_is_institute_book_idx;

ALTER TABLE books DROP COLUMN IF EXISTS is_institute_book;
