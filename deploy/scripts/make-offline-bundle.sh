#!/bin/bash
# Сборка офлайн-bundle для деплоя в институте.
# Запускается на машине С ИНТЕРНЕТОМ.
# На выходе: ../library-offline-bundle-YYYY-MM-DD.tar.gz
#
# Содержимое bundle:
#   images.tar              — все docker-образы (api, frontend, postgres, nginx)
#   docker-compose.prod.yml
#   .env.example
#   nginx/                  — nginx-конфиги
#   postgres/init.sql
#   scripts/                — install.sh, uninstall.sh, backup.sh, restore.sh, …
#   docker-engine/          — .deb-пакеты Docker Engine + compose plugin (Ubuntu 22.04/24.04)
#   docs/INSTALL_OFFLINE.md

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
DEPLOY="$ROOT/deploy"
TAG="${TAG:-latest}"
REGISTRY="${REGISTRY:-localhost}"
DATE="$(date +%Y-%m-%d)"
BUNDLE_NAME="library-offline-bundle-$DATE"
WORK_DIR="$(mktemp -d -t library-bundle-XXXXXX)"
OUT_DIR="$WORK_DIR/$BUNDLE_NAME"
OUT_ARCHIVE="$ROOT/$BUNDLE_NAME.tar.gz"

DOCKER_DEB_DIR="$DEPLOY/docker-engine"
DOCKER_VERSION="${DOCKER_VERSION:-5:27.3.1-1~ubuntu.22.04~jammy}"

cleanup() { rm -rf "$WORK_DIR"; }
trap cleanup EXIT

log() { printf '\033[1;34m[bundle]\033[0m %s\n' "$*"; }

log "1/5 сборка docker-образов (tag: $TAG)"
docker build --pull -t "$REGISTRY/library-api:$TAG" "$ROOT/api"
docker build --pull -t "$REGISTRY/library-frontend:$TAG" "$ROOT/frontend"
docker pull postgres:16-alpine
docker pull nginx:1.27-alpine

log "2/5 экспорт образов → images.tar"
mkdir -p "$OUT_DIR"
docker save \
  "$REGISTRY/library-api:$TAG" \
  "$REGISTRY/library-frontend:$TAG" \
  postgres:16-alpine \
  nginx:1.27-alpine \
  -o "$OUT_DIR/images.tar"

log "3/5 копирование конфигов и скриптов"
cp "$DEPLOY/docker-compose.prod.yml" "$OUT_DIR/"
cp "$DEPLOY/.env.example"            "$OUT_DIR/"
cp -r "$DEPLOY/nginx"                "$OUT_DIR/"
cp -r "$DEPLOY/postgres"             "$OUT_DIR/"
cp -r "$DEPLOY/scripts"              "$OUT_DIR/"
chmod +x "$OUT_DIR/scripts/"*.sh
if [ -d "$ROOT/docs" ]; then
  mkdir -p "$OUT_DIR/docs"
  cp "$ROOT/docs/INSTALL_OFFLINE.md" "$OUT_DIR/docs/" 2>/dev/null || true
fi

log "4/5 скачивание .deb пакетов Docker Engine (для офлайн-установки)"
if [ ! -d "$DOCKER_DEB_DIR" ] || [ -z "$(ls -A "$DOCKER_DEB_DIR" 2>/dev/null)" ]; then
  if command -v apt-get >/dev/null 2>&1; then
    mkdir -p "$DOCKER_DEB_DIR"
    pushd "$DOCKER_DEB_DIR" >/dev/null
    # Тянем .deb на хост-машине через apt-download (или вручную с download.docker.com)
    apt-get download \
      docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin \
      2>/dev/null || true
    popd >/dev/null
  fi
  if [ -z "$(ls -A "$DOCKER_DEB_DIR" 2>/dev/null)" ]; then
    log "  (пропущено: нет apt-get или не Ubuntu/Debian.)"
    log "  Положи .deb пакеты Docker Engine вручную в $DOCKER_DEB_DIR и пересобери."
    log "  Скачать: https://download.docker.com/linux/ubuntu/dists/jammy/pool/stable/amd64/"
  fi
fi
cp -r "$DOCKER_DEB_DIR" "$OUT_DIR/" 2>/dev/null || mkdir -p "$OUT_DIR/docker-engine"

log "5/5 упаковка → $OUT_ARCHIVE"
tar -czf "$OUT_ARCHIVE" -C "$WORK_DIR" "$BUNDLE_NAME"

SIZE="$(du -h "$OUT_ARCHIVE" | awk '{print $1}')"
log "готово: $OUT_ARCHIVE ($SIZE)"
log ""
log "Следующие шаги:"
log "  1. Перенеси архив на сервер института (USB/scp)."
log "  2. На сервере: tar xzf $BUNDLE_NAME.tar.gz && cd $BUNDLE_NAME"
log "  3. Запусти: sudo ./scripts/install.sh"
