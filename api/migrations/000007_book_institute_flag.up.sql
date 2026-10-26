ALTER TABLE books ADD COLUMN is_institute_book BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX books_is_institute_book_idx ON books(is_institute_book);
