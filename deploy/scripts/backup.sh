#!/bin/sh
# Ночной бэкап PostgreSQL.
# Запускается из контейнера backup (см. docker-compose.prod.yml).
# Хранит:
#   /backups/daily/library-YYYY-MM-DD.sql.gz   — 7 копий
#   /backups/weekly/library-YYYY-WW.sql.gz     — 4 копии (по воскресеньям)

set -eu

: "${POSTGRES_USER:?POSTGRES_USER is required}"
: "${POSTGRES_PASSWORD:?POSTGRES_PASSWORD is required}"
: "${POSTGRES_DB:?POSTGRES_DB is required}"
: "${PGHOST:=postgres}"

KEEP_DAILY="${BACKUP_KEEP_DAILY:-7}"
KEEP_WEEKLY="${BACKUP_KEEP_WEEKLY:-4}"

DAILY_DIR=/backups/daily
WEEKLY_DIR=/backups/weekly
mkdir -p "$DAILY_DIR" "$WEEKLY_DIR"

TODAY="$(date +%Y-%m-%d)"
WEEK="$(date +%Y-%V)"
DAILY_FILE="$DAILY_DIR/library-$TODAY.sql.gz"
WEEKLY_FILE="$WEEKLY_DIR/library-$WEEK.sql.gz"

echo "[backup] $(date -Iseconds) starting pg_dump"

export PGPASSWORD="$POSTGRES_PASSWORD"
pg_dump --host="$PGHOST" --username="$POSTGRES_USER" --dbname="$POSTGRES_DB" \
        --format=plain --no-owner --no-privileges --clean --if-exists \
  | gzip -9 > "$DAILY_FILE.tmp"
mv "$DAILY_FILE.tmp" "$DAILY_FILE"

# По воскресеньям также делаем еженедельную копию (hard link, чтобы не дублировать на диске)
if [ "$(date +%u)" = "7" ]; then
  ln -f "$DAILY_FILE" "$WEEKLY_FILE"
fi

# Ротация: оставляем только последние N
ls -1t "$DAILY_DIR"/library-*.sql.gz 2>/dev/null | tail -n +"$((KEEP_DAILY + 1))" | xargs -r rm -f
ls -1t "$WEEKLY_DIR"/library-*.sql.gz 2>/dev/null | tail -n +"$((KEEP_WEEKLY + 1))" | xargs -r rm -f

SIZE="$(du -h "$DAILY_FILE" | awk '{print $1}')"
echo "[backup] $(date -Iseconds) done: $DAILY_FILE ($SIZE)"
