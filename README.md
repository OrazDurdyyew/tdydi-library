# Электронная библиотека

Универсальная платформа электронной библиотеки для ВУЗов.

См. [ТЗ.md](ТЗ.md) — техническое задание, [CLAUDE.md](CLAUDE.md) — контекст проекта.

## Стек

- **Frontend**: Nuxt 4 + TypeScript + TailwindCSS + shadcn-vue + i18n (TM/RU/EN)
- **Backend**: Go + Fiber + GORM
- **БД**: PostgreSQL 16 (FTS + pg_trgm)
- **Деплой**: Docker Compose + приватный Docker Registry

## Быстрый старт (разработка)

### Требования
- Docker и Docker Compose
- (опционально для локальной разработки без Docker) Node.js 22+, Go 1.23+

### Запуск
```bash
make dev          # одна команда: поднимает всё и показывает логи
# или вручную:
cp .env.example .env
cd deploy
docker compose -f docker-compose.dev.yml up -d
```

Полезные команды: `make help`.

После старта:
- Frontend: http://localhost:2026
- API:      http://localhost:2020
- API health: http://localhost:2020/api/health
- PostgreSQL: localhost:5433 (library / library) — внутри docker-сети 5432

### Остановить
```bash
cd deploy
docker compose -f docker-compose.dev.yml down
```

### Сбросить БД (удалить volumes)
```bash
docker compose -f docker-compose.dev.yml down -v
```

## Структура проекта

```
/api          # Go (Fiber) backend
  /cmd/server # точка входа
  /internal   # модели, handlers, services, storage
  /migrations # golang-migrate
/frontend     # Nuxt 4 (публичный сайт + /admin/*)
  /pages
  /layouts
  /middleware
  /components
  /locales
/deploy       # docker-compose, nginx, init-скрипты
/docs         # документация развёртывания
.claude/skills # вспомогательные команды для разработки
```

## Этапы разработки

См. [ТЗ.md](ТЗ.md) раздел 9. Сейчас — старт MVP.

## Развёртывание для нового ВУЗа

См. [docs/deploy-new-vuz.md](docs/deploy-new-vuz.md).
