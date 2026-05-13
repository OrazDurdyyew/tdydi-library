#!/bin/bash
# Восстановление БД из бэкапа.
# Запускать на хосте (НЕ внутри контейнера).
#
# Использование:
#   ./restore.sh                            — интерактивный выбор файла
#   ./restore.sh /var/library/backups/daily/library-2026-05-13.sql.gz
#
# ВНИМАНИЕ: дроп всех существующих таблиц перед восстановлением.

set -euo pipefail

STORAGE_ROOT="${STORAGE_ROOT:-/var/library}"
BACKUPS_DIR="$STORAGE_ROOT/backups"
COMPOSE_FILE="$(dirname "$0")/../docker-compose.prod.yml"
ENV_FILE="$(dirname "$0")/../.env"

if [ ! -f "$ENV_FILE" ]; then
  echo "Не найден $ENV_FILE. Скопируй .env.example в .env и заполни."
  exit 1
fi

# shellcheck disable=SC1090
set -a; . "$ENV_FILE"; set +a

# === Выбор файла ===
FILE="${1:-}"
if [ -z "$FILE" ]; then
  echo "Доступные бэкапы:"
  echo "--- daily ---"
  ls -1t "$BACKUPS_DIR/daily/" 2>/dev/null || echo "  (нет)"
  echo "--- weekly ---"
  ls -1t "$BACKUPS_DIR/weekly/" 2>/dev/null || echo "  (нет)"
  echo
  read -rp "Введи полный путь к .sql.gz: " FILE
fi

if [ ! -f "$FILE" ]; then
  echo "Файл не найден: $FILE"
  exit 1
fi

echo
echo "Будет восстановлена БД '$POSTGRES_DB' из:"
echo "  $FILE"
echo
read -rp "Все текущие данные будут удалены. Продолжить? (yes/NO): " CONFIRM
[ "$CONFIRM" = "yes" ] || { echo "Отменено."; exit 0; }

echo "[restore] остановка api/frontend (postgres остаётся)…"
docker compose -f "$COMPOSE_FILE" stop api frontend backup

echo "[restore] восстановление…"
gunzip -c "$FILE" | docker compose -f "$COMPOSE_FILE" exec -T postgres \
  psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1

echo "[restore] запуск сервисов…"
docker compose -f "$COMPOSE_FILE" up -d

echo "[restore] готово."
