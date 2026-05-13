---
name: run-migrations
description: Применить миграции к локальной БД через golang-migrate. Триггер — когда пользователь говорит "накати миграции", "обнови БД", "примени миграции", "migrate up". Запускает golang-migrate в нужном режиме.
---

# run-migrations

Применяет миграции БД через `golang-migrate`.

## Когда использовать

- "накати миграции"
- "обнови БД до последней версии"
- "примени миграции"
- после создания новой миграции через `/new-migration`

## Команды

### Применить все pending миграции
```sh
migrate -path api/migrations -database "$DATABASE_URL" up
```

### Откатить последнюю миграцию
```sh
migrate -path api/migrations -database "$DATABASE_URL" down 1
```

### Откатить N миграций
```sh
migrate -path api/migrations -database "$DATABASE_URL" down N
```

### Применить N следующих
```sh
migrate -path api/migrations -database "$DATABASE_URL" up N
```

### Текущая версия
```sh
migrate -path api/migrations -database "$DATABASE_URL" version
```

### Force версию (для битого dirty-state)
```sh
migrate -path api/migrations -database "$DATABASE_URL" force VERSION
```

## DATABASE_URL для разработки

```
postgres://library:library@localhost:5432/library?sslmode=disable
```

(берётся из `.env` или docker-compose окружения)

## Когда использовать `force`

Если миграция упала в середине → состояние `dirty`. Тогда:
1. Вручную исправить БД (или дописать `up.sql` чтобы он был idempotent через `IF NOT EXISTS`)
2. `migrate force VERSION` — где VERSION = последняя успешная (предыдущая)
3. Снова `migrate up`

**Никогда не использовать `force` на production без бэкапа.**

## Если migrate не установлен

```sh
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Или использовать Docker:
```sh
docker run --rm -v $(pwd)/api/migrations:/migrations migrate/migrate \
  -path=/migrations -database "$DATABASE_URL" up
```
