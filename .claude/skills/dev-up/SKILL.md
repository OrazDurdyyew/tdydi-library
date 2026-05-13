---
name: dev-up
description: Запустить локальное окружение разработки через docker-compose. Триггер — когда пользователь говорит "подними проект", "запусти dev", "старт окружения", "поехали". Запускает PostgreSQL, бэк, фронт, накатывает миграции, делает seed (опционально).
---

# dev-up

Запускает локальное окружение для разработки.

## Когда использовать

- "подними проект"
- "запусти dev окружение"
- "старт"
- в начале новой сессии разработки

## Сервисы в docker-compose.yml (dev)

- `postgres` — PostgreSQL 16 с расширениями `pg_trgm`, `unaccent` (через init.sql)
- `api` — Go-бэкенд (hot-reload через `air`)
- `frontend` — Nuxt dev-сервер (`npm run dev`)
- (опционально) `redis` — если используется кэш

## Команды

### Первый запуск
```sh
docker compose up -d postgres
# подождать пока postgres стартанёт
# применить миграции
/run-migrations
# опционально — наполнить тестовыми данными
/seed-db
# запустить api и frontend
docker compose up -d api frontend
```

### Повседневный запуск
```sh
docker compose up -d
```

### Логи
```sh
docker compose logs -f api          # логи бэка
docker compose logs -f frontend     # логи фронта
docker compose logs -f postgres     # логи БД
```

### Остановить всё
```sh
docker compose down
```

### Сбросить БД и начать заново
```sh
docker compose down -v              # удаляет volumes
docker compose up -d postgres
/run-migrations
/seed-db
```

## Порты по умолчанию (dev)

- Frontend: `http://localhost:2026`
- API: `http://localhost:2020`
- PostgreSQL: `localhost:5432`

## .env для разработки

Должен лежать в корне проекта (НЕ в git):
```
DATABASE_URL=postgres://library:library@postgres:5432/library?sslmode=disable
JWT_SECRET=dev-secret-not-for-production
STORAGE_PATH=/var/library
NUXT_PUBLIC_API_BASE=http://localhost:2020/api
```

Шаблон — `.env.example` в репо.

## Если что-то не работает

- Проверить логи через `docker compose logs -f`
- Проверить что миграции применены: `migrate version`
- Если "dirty" state в миграциях — см. `/run-migrations` skill
- Если порт занят — проверить `lsof -i :2026` / `lsof -i :2020`
