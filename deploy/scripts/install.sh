#!/bin/bash
# Установка/обновление платформы на сервере института.
# Запускать ОТ ROOT (или через sudo).
#
# Использование:
#   sudo ./install.sh           — полная установка (первая)
#   sudo ./install.sh --update  — обновление (только подменить образы и перезапустить)

set -euo pipefail

# ===== Цвета =====
C_BLUE='\033[1;34m'; C_GREEN='\033[1;32m'; C_RED='\033[1;31m'; C_YELLOW='\033[1;33m'; C_RST='\033[0m'
log()  { printf "${C_BLUE}[install]${C_RST} %s\n" "$*"; }
ok()   { printf "${C_GREEN}[install]${C_RST} %s\n" "$*"; }
warn() { printf "${C_YELLOW}[install]${C_RST} %s\n" "$*"; }
die()  { printf "${C_RED}[install]${C_RST} %s\n" "$*" >&2; exit 1; }

# ===== Проверка прав =====
if [ "$(id -u)" -ne 0 ]; then
  die "Запусти через sudo: sudo $0"
fi

UPDATE_MODE=false
[ "${1:-}" = "--update" ] && UPDATE_MODE=true

BUNDLE_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$BUNDLE_DIR"

# ===== 1. Установка Docker (если ещё нет) =====
if ! command -v docker >/dev/null 2>&1; then
  log "Docker не найден, ставлю из локальных .deb пакетов…"
  if [ -d "$BUNDLE_DIR/docker-engine" ] && [ -n "$(ls -A "$BUNDLE_DIR/docker-engine" 2>/dev/null)" ]; then
    dpkg -i "$BUNDLE_DIR/docker-engine/"*.deb || apt-get install -f -y
    systemctl enable --now docker
    ok "Docker установлен."
  else
    die "Нет пакетов в $BUNDLE_DIR/docker-engine/. Установи Docker вручную или пересобери bundle с .deb."
  fi
else
  ok "Docker уже установлен ($(docker --version))."
fi

if ! docker compose version >/dev/null 2>&1; then
  die "Docker Compose plugin не найден. Установи docker-compose-plugin."
fi

# ===== 2. .env =====
if [ ! -f "$BUNDLE_DIR/.env" ]; then
  log "Создаю .env из .env.example…"
  cp "$BUNDLE_DIR/.env.example" "$BUNDLE_DIR/.env"

  # Автогенерация секретов
  PG_PASS="$(openssl rand -base64 32 | tr -d '/+=' | head -c 32)"
  JWT_SEC="$(openssl rand -base64 48 | tr -d '/+=' | head -c 48)"
  sed -i "s|POSTGRES_PASSWORD=.*|POSTGRES_PASSWORD=$PG_PASS|" "$BUNDLE_DIR/.env"
  sed -i "s|JWT_SECRET=.*|JWT_SECRET=$JWT_SEC|" "$BUNDLE_DIR/.env"
  chmod 600 "$BUNDLE_DIR/.env"
  ok ".env создан, секреты сгенерированы автоматически."
  warn "Сохрани содержимое $BUNDLE_DIR/.env в надёжное место (это пароли БД и подпись JWT)."
else
  ok ".env уже есть, не трогаю."
fi

# ===== 2b. Синхронизация TAG с версией bundle =====
# Образы тегируются версией (не latest), а compose берёт образ по ${TAG}.
# При обновлении переносится СТАРЫЙ .env со старым TAG — без этой синхронизации
# compose поднял бы старые образы и обновление не применилось бы.
if [ -f "$BUNDLE_DIR/VERSION" ]; then
  BUNDLE_VERSION="$(tr -d '[:space:]' < "$BUNDLE_DIR/VERSION")"
  if [ -n "$BUNDLE_VERSION" ]; then
    if grep -q '^TAG=' "$BUNDLE_DIR/.env"; then
      sed -i "s|^TAG=.*|TAG=$BUNDLE_VERSION|" "$BUNDLE_DIR/.env"
    else
      printf 'TAG=%s\n' "$BUNDLE_VERSION" >> "$BUNDLE_DIR/.env"
    fi
    ok "Версия bundle: $BUNDLE_VERSION (TAG в .env синхронизирован)."
  fi
fi

# Подхватываем переменные
set -a; . "$BUNDLE_DIR/.env"; set +a
STORAGE_ROOT="${STORAGE_ROOT:-/var/library}"

# ===== 3. Директории хранилища =====
log "Создаю директории в $STORAGE_ROOT…"
mkdir -p \
  "$STORAGE_ROOT/pg-data" \
  "$STORAGE_ROOT/files/books" \
  "$STORAGE_ROOT/files/periodicals" \
  "$STORAGE_ROOT/files/branding" \
  "$STORAGE_ROOT/backups/daily" \
  "$STORAGE_ROOT/backups/weekly" \
  "$STORAGE_ROOT/ssl"

# PostgreSQL в alpine-образе работает под UID 70 (postgres)
chown -R 70:70 "$STORAGE_ROOT/pg-data"
chmod 700 "$STORAGE_ROOT/pg-data"
# Nginx внутри nginx:alpine — UID 101
chown -R 101:101 "$STORAGE_ROOT/ssl" 2>/dev/null || true
# Файлы книг доступны и api (UID 1000 из Dockerfile), и nginx (101) на чтение
chown -R 1000:1000 "$STORAGE_ROOT/files"
chmod -R 755 "$STORAGE_ROOT/files"
chown -R 0:0 "$STORAGE_ROOT/backups"
ok "Директории готовы."

# ===== 4. Загрузка образов =====
if [ -f "$BUNDLE_DIR/images.tar" ]; then
  log "Загрузка docker-образов (это может занять минуту)…"
  docker load -i "$BUNDLE_DIR/images.tar"
  ok "Образы загружены."
else
  warn "images.tar не найден — предполагаю, что образы уже в локальном реестре."
fi

# ===== 5. Запуск =====
COMPOSE="docker compose -f $BUNDLE_DIR/docker-compose.prod.yml --env-file $BUNDLE_DIR/.env"

if $UPDATE_MODE; then
  log "Режим обновления: перезапускаю api/frontend…"
  $COMPOSE up -d --force-recreate api frontend nginx
else
  log "Запуск всего стэка…"
  $COMPOSE up -d
fi

log "Жду готовности сервисов…"
sleep 5
$COMPOSE ps

# ===== 6. Создание первого суперадмина =====
if ! $UPDATE_MODE; then
  echo
  read -rp "Создать первого суперадмина сейчас? (y/N): " CREATE_ADMIN
  if [ "$CREATE_ADMIN" = "y" ] || [ "$CREATE_ADMIN" = "Y" ]; then
    read -rp "Логин: " ADMIN_USER
    read -rsp "Пароль: " ADMIN_PASS; echo
    $COMPOSE exec -T api /app/createadmin \
      -username="$ADMIN_USER" -password="$ADMIN_PASS" -role=superadmin \
      || warn "Не получилось через CLI — создай вручную (см. INSTALL_OFFLINE.md)."
  fi
fi

echo
ok "Готово. Установленная версия: ${TAG:-неизвестно}"
ok "Проверка версии:  curl -s http://127.0.0.1${HTTP_PORT:+:$HTTP_PORT}/api/health"
ok "Сайт:    http://$(hostname -I | awk '{print $1}')${HTTP_PORT:+:$HTTP_PORT}/"
ok "Админка: http://$(hostname -I | awk '{print $1}')${HTTP_PORT:+:$HTTP_PORT}/admin/login"
echo
log "Полезные команды:"
log "  Статус:        $COMPOSE ps"
log "  Логи api:      $COMPOSE logs -f api"
log "  Бэкап вручную: $COMPOSE exec backup sh /usr/local/bin/backup.sh"
log "  Восстановить:  ./scripts/restore.sh"
log "  Удалить:       sudo ./scripts/uninstall.sh"
