---
name: new-migration
description: Создать новую миграцию БД для golang-migrate. Триггер — когда пользователь просит добавить таблицу, изменить схему, добавить колонку, индекс, расширение в PostgreSQL. Создаёт пару файлов up/down с правильной нумерацией.
---

# new-migration

Создаёт пару миграционных файлов в `/api/migrations/` для PostgreSQL через golang-migrate.

## Когда использовать

- "добавь таблицу X"
- "нужна миграция для ..."
- "измени схему ..."
- "добавь колонку / индекс / расширение"

## Что сделать

1. Прочитай существующие миграции в `/api/migrations/` чтобы понять текущий номер.
2. Определи следующий номер: `{N+1}` с padding до 6 цифр (`000001_init.up.sql`).
3. Спроси у пользователя короткое имя миграции (например, `add_books_table`), если не очевидно.
4. Создай два файла:
   - `{NNNNNN}_{name}.up.sql` — изменения
   - `{NNNNNN}_{name}.down.sql` — обратные изменения (DROP, ALTER ... DROP COLUMN и т.д.)
5. В `up.sql` для новых таблиц добавь:
   - Все нужные индексы (PK, FK, уникальные, GIN для tsvector, B-tree для часто фильтруемых полей)
   - Триггеры для `search_vector` если поле есть
   - `CHECK` constraints где применимо (например, для 2-уровневой иерархии `sections`)
6. В `down.sql` всегда обратные операции — корректные `DROP TABLE IF EXISTS`, `DROP INDEX`, и т.д.

## Соглашения проекта

- Названия таблиц — `snake_case` множественное число (`books`, `book_sections`)
- Названия колонок — `snake_case` единственное (`book_id`, `created_at`)
- Все таблицы имеют `id BIGSERIAL PRIMARY KEY` и `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()` где применимо
- M:N таблицы — без `id`, составной PK из двух FK
- Поля на 3 языках: `name_tm`, `name_ru`, `name_en` (всегда NOT NULL)
- Текстовые поля в одном языке книги: `title`, `annotation`, `authors`
- FK с `ON DELETE` указывать явно (`CASCADE` для зависимых, `RESTRICT` для критичных)
- Все timestamps — `TIMESTAMPTZ`, не `TIMESTAMP`

## После создания

- Напомнить пользователю запустить миграции через `/run-migrations` или `migrate -path api/migrations -database "$DB_URL" up`
- Если изменилась модель — обновить соответствующий Go-файл в `/api/internal/models/`
