---
name: new-handler
description: Создать CRUD handler в Go для новой сущности. Триггер — когда пользователь просит добавить API-эндпоинт для модели (книги, разделы, периодика, новости и т.д.). Создаёт модель, handler с CRUD-методами, регистрирует роуты, добавляет проверку прав.
---

# new-handler

Создаёт полный CRUD для новой сущности в Go (Fiber + GORM).

## Когда использовать

- "добавь API для X"
- "создай эндпоинты для книг / разделов / периодики"
- "нужен handler для ..."

## Что создать

1. **Модель** в `/api/internal/models/{entity}.go`:
   - GORM-структура с тегами
   - Связи через `gorm:"foreignKey:..."` или `many2many:...`
   - Методы валидации
2. **Repository** в `/api/internal/storage/{entity}_repo.go`:
   - `List(ctx, filter, page, limit)`
   - `Get(ctx, id)`
   - `Create(ctx, entity)`
   - `Update(ctx, id, updates)`
   - `Delete(ctx, id)`
   - Для сущностей с поиском — фильтр по `search_vector` через FTS
3. **Handler** в `/api/internal/handlers/{entity}.go`:
   - `List` (GET /api/...): пагинация, фильтры, сортировка
   - `Get` (GET /api/.../:id)
   - `Create` (POST /api/admin/...) — требует auth + permission check
   - `Update` (PUT /api/admin/.../:id) — auth + permission
   - `Delete` (DELETE /api/admin/.../:id) — auth + permission
4. **DTO** для request/response с тегами `validate` (валидация через go-playground/validator)
5. **Роуты** в `/api/cmd/main.go` (или `routes.go`):
   - Публичные `/api/...` — без middleware
   - Админские `/api/admin/...` — `auth` middleware + опционально `section_permission` middleware

## Проверка прав (важно)

Для эндпоинтов работы с книгами и разделами:
- Все админские эндпоинты — JWT middleware
- Для книг: после auth — проверка `admin_section_permissions` (или `all_sections=true`).
  - На `Create`: проверить что все `section_ids` из тела запроса есть в разрешённых
  - На `Update`/`Delete`: получить текущую книгу, проверить что **все** её разделы — в разрешённых
- Все остальные сущности (sections, periodicals, news, faq, pages, settings) — только `superadmin`

## Соглашения

- Имена методов: `ListBooks`, `GetBook`, `CreateBook`, `UpdateBook`, `DeleteBook`
- Response: `{ "data": [...], "meta": { "total": N, "page": P, "limit": L } }` для списков
- Ошибки: `{ "error": "message", "code": "BOOK_NOT_FOUND" }` с HTTP-статусом
- Логирование: при критичных ошибках через `slog`

## Шаблоны (когда появятся в проекте)

Поискать в `/api/internal/handlers/` уже существующий handler и использовать его как шаблон.
