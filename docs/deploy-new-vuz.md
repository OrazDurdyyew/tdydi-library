# Развёртывание платформы для нового ВУЗа

## Требования к серверу

- Linux (рекомендуется Ubuntu 22.04 LTS или новее)
- Минимум 4 ГБ RAM
- SSD от 200 ГБ (под файлы книг)
- Docker 24+ и Docker Compose v2
- Домен с A-записью на сервер
- SSL-сертификат (Let's Encrypt через certbot)

## Шаги

### 1. Доступ к приватному Docker Registry

Разработчик предоставляет учётные данные. Логин:

```bash
docker login ghcr.io
# или другой registry: harbor.example.com, docker.io
```

### 2. Подготовка рабочей директории

```bash
sudo mkdir -p /opt/library
cd /opt/library
```

### 3. Скачать `docker-compose.yml` и `.env.example`

Файлы предоставляет разработчик. Положить в `/opt/library/`.

### 4. Настроить `.env`

```bash
cp .env.example .env
nano .env
```

Обязательно изменить:
- `POSTGRES_PASSWORD` — сильный пароль для БД
- `JWT_SECRET` — длинная случайная строка (минимум 32 символа)
- `REGISTRY` и `TAG` — координаты образов
- `ALLOWED_ORIGINS` — домен ВУЗа

### 5. Скачать образы и запустить

```bash
docker compose pull
docker compose up -d
```

### 6. Применить миграции (первый запуск)

```bash
docker compose exec api migrate -path /app/migrations -database "$DATABASE_URL" up
```

### 7. Создать первого суперадмина

```bash
docker compose exec api /app/library-api create-admin \
  --username=admin \
  --password=<сильный-пароль>
```

### 8. Настроить SSL и reverse proxy

Вариант А — Caddy (рекомендуется, авто-SSL):

```caddy
library.vuz.tm {
    reverse_proxy localhost:2026
    reverse_proxy /api/* localhost:2020
    reverse_proxy /protected/* localhost:2020
}
```

Вариант Б — Nginx + certbot: настроить как обычно.

### 9. Первичная настройка через админку

Открыть `https://library.vuz.tm/admin/login` → войти → пройти мастер:
- Название института на 3 языках
- Логотип, фавикон
- Цвета палитры
- Контакты

## Обновление до новой версии

```bash
cd /opt/library
docker compose pull
docker compose up -d
docker compose exec api migrate -path /app/migrations -database "$DATABASE_URL" up
```

## Бэкап

Создать `/opt/library/backup.sh`:

```bash
#!/bin/bash
set -e
DATE=$(date +%Y%m%d)
BACKUP_DIR=/backups/library

mkdir -p $BACKUP_DIR

# БД
docker compose -f /opt/library/docker-compose.yml exec -T postgres \
  pg_dump -U library library | gzip > $BACKUP_DIR/db-$DATE.sql.gz

# Файлы (книги, обложки, брендинг)
rsync -a /var/lib/docker/volumes/library_library-storage/_data/ \
  $BACKUP_DIR/files-$DATE/

# Удалить старые: оставить 7 дневных + 4 недельных
find $BACKUP_DIR -name 'db-*.sql.gz' -mtime +7 ! -name '*-monday-*' -delete
find $BACKUP_DIR -name 'files-*' -mtime +7 ! -name '*-monday-*' -delete
```

```bash
chmod +x /opt/library/backup.sh
echo "0 3 * * * /opt/library/backup.sh" | sudo crontab -
```

## Восстановление

```bash
# БД
gunzip < db-20260101.sql.gz | docker compose exec -T postgres psql -U library library

# Файлы
rsync -a files-20260101/ /var/lib/docker/volumes/library_library-storage/_data/
```

## Мониторинг

Проверить здоровье:
- `curl https://library.vuz.tm/api/health` → `{"status":"ok"}`
- `docker compose ps` — все контейнеры `running`/`healthy`

## Поддержка

При проблемах:
1. `docker compose logs -f api` / `frontend` / `postgres`
2. Связаться с разработчиком, предоставив логи
