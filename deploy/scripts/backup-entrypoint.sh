#!/bin/sh
# Entrypoint бэкап-контейнера: запускает backup.sh по cron-расписанию.
# По умолчанию: каждый день в 03:00.

set -eu

BACKUP_CRON="${BACKUP_CRON:-0 3 * * *}"

mkdir -p /var/spool/cron/crontabs

# Прокидываем env-переменные внутрь cron (он их по умолчанию не видит)
{
  echo "POSTGRES_USER=$POSTGRES_USER"
  echo "POSTGRES_PASSWORD=$POSTGRES_PASSWORD"
  echo "POSTGRES_DB=$POSTGRES_DB"
  echo "PGHOST=${PGHOST:-postgres}"
  echo "BACKUP_KEEP_DAILY=${BACKUP_KEEP_DAILY:-7}"
  echo "BACKUP_KEEP_WEEKLY=${BACKUP_KEEP_WEEKLY:-4}"
  echo ""
  echo "$BACKUP_CRON sh /usr/local/bin/backup.sh >> /var/log/backup.log 2>&1"
} > /var/spool/cron/crontabs/root

touch /var/log/backup.log

echo "[backup] entrypoint: scheduled cron '$BACKUP_CRON'"
echo "[backup] running first backup at startup to verify configuration..."
sh /usr/local/bin/backup.sh || echo "[backup] initial run failed (БД может быть пустой при первом запуске — это нормально)"

# Запускаем crond на foreground + хвост лога
crond -f -l 8 &
CROND_PID=$!

tail -F /var/log/backup.log &
TAIL_PID=$!

trap "kill $CROND_PID $TAIL_PID 2>/dev/null; exit 0" TERM INT
wait $CROND_PID
