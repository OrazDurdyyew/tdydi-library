#!/bin/bash
# Удаление платформы с сервера.
# По умолчанию: только останавливает контейнеры (данные остаются).
# С --purge: удаляет ВСЕ данные и образы.

set -euo pipefail

C_BLUE='\033[1;34m'; C_RED='\033[1;31m'; C_YELLOW='\033[1;33m'; C_RST='\033[0m'
log()  { printf "${C_BLUE}[uninstall]${C_RST} %s\n" "$*"; }
warn() { printf "${C_YELLOW}[uninstall]${C_RST} %s\n" "$*"; }
die()  { printf "${C_RED}[uninstall]${C_RST} %s\n" "$*" >&2; exit 1; }

[ "$(id -u)" -eq 0 ] || die "Запусти через sudo."

BUNDLE_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PURGE=false
[ "${1:-}" = "--purge" ] && PURGE=true

COMPOSE="docker compose -f $BUNDLE_DIR/docker-compose.prod.yml"
[ -f "$BUNDLE_DIR/.env" ] && COMPOSE="$COMPOSE --env-file $BUNDLE_DIR/.env"

log "Останавливаю контейнеры…"
$COMPOSE down

if ! $PURGE; then
  log "Контейнеры удалены. Данные сохранены в \$STORAGE_ROOT (по умолчанию /var/library/)."
  log "Чтобы стереть ВСЁ (включая БД, файлы и бэкапы): sudo $0 --purge"
  exit 0
fi

warn "═══════════════════════════════════════════════════════════"
warn "  --purge удалит ВСЕ данные платформы (БД, файлы, бэкапы)."
warn "═══════════════════════════════════════════════════════════"
read -rp "Введи 'DELETE EVERYTHING' для подтверждения: " CONFIRM
[ "$CONFIRM" = "DELETE EVERYTHING" ] || { log "Отменено."; exit 0; }

# Подхватываем STORAGE_ROOT
[ -f "$BUNDLE_DIR/.env" ] && { set -a; . "$BUNDLE_DIR/.env"; set +a; }
STORAGE_ROOT="${STORAGE_ROOT:-/var/library}"

log "Удаление данных $STORAGE_ROOT…"
rm -rf "$STORAGE_ROOT"

log "Удаление docker-образов…"
docker image rm \
  "${REGISTRY:-localhost}/library-api:${TAG:-latest}" \
  "${REGISTRY:-localhost}/library-frontend:${TAG:-latest}" \
  postgres:16-alpine \
  nginx:1.27-alpine \
  2>/dev/null || true

log "Готово. .env и bundle оставлены — удали вручную если нужно."
