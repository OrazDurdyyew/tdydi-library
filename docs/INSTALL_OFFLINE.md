# Установка платформы «Электронная библиотека» на офлайн-сервер

Эта инструкция для администратора ВУЗа, разворачивающего платформу на локальном сервере **без доступа к интернету**.

---

## Системные требования

- **ОС**: Ubuntu Server 22.04 LTS или 24.04 LTS (рекомендуется), Debian 12.
- **CPU**: 2+ ядра.
- **RAM**: минимум 4 ГБ (рекомендуется 8 ГБ).
- **Диск**: SSD от 200 ГБ. Данные хранятся в `/var/library/` — выдели этой директории отдельный раздел, если возможно.
- **Сеть**: статический IP в локальной сети ВУЗа. Опционально — домен (например, `library.institute.local`).

Сервер может быть полностью офлайн — интернет нужен только разработчику (один раз) для сборки bundle.

---

## Что входит в bundle

Архив `library-offline-bundle-YYYY-MM-DD.tar.gz` содержит:

```
library-offline-bundle-YYYY-MM-DD/
├── images.tar                  # docker-образы (api, frontend, postgres, nginx) — ~700 МБ
├── docker-compose.prod.yml     # описание стэка
├── .env.example                # шаблон конфигурации
├── nginx/                      # конфиги reverse proxy
├── postgres/init.sql           # расширения БД (pg_trgm, unaccent)
├── docker-engine/              # .deb пакеты Docker Engine (~150 МБ)
└── scripts/
    ├── install.sh              # установщик
    ├── uninstall.sh            # удаление
    ├── backup.sh               # скрипт ночного бэкапа
    ├── restore.sh              # восстановление из бэкапа
    └── backup-entrypoint.sh    # entrypoint бэкап-сайдкара
```

Общий размер архива — ~1 ГБ.

---

## Шаг 1. Перенос архива на сервер

С флешки или через `scp` из локальной сети:

```bash
# Скопировать архив в домашнюю директорию администратора
scp library-offline-bundle-2026-05-13.tar.gz admin@library-server:~/
```

---

## Шаг 2. Распаковка

На сервере:

```bash
cd ~
tar xzf library-offline-bundle-2026-05-13.tar.gz
cd library-offline-bundle-2026-05-13
```

---

## Шаг 3. Запуск установщика

```bash
sudo ./scripts/install.sh
```

Установщик автоматически:

1. Проверит наличие Docker. Если его нет — установит из `docker-engine/*.deb`.
2. Создаст файл `.env` со случайно сгенерированными паролями БД и подписи JWT.
3. Создаст директории в `/var/library/`:
   - `pg-data/` — данные PostgreSQL
   - `files/books/` — файлы книг (PDF/EPUB)
   - `files/periodicals/` — выпуски периодики
   - `files/branding/` — логотип, фавикон
   - `backups/daily/` — ежедневные бэкапы (7 копий)
   - `backups/weekly/` — еженедельные (4 копии)
   - `ssl/` — место под SSL-сертификаты
4. Загрузит docker-образы из `images.tar`.
5. Запустит весь стек.
6. Спросит, создать ли первого суперадмина.

После установки сайт доступен на `http://<ip-сервера>/`, админка — на `http://<ip-сервера>/admin/login`.

> **Важно**: содержимое `/root/library-offline-bundle-*/​.env` нужно сохранить в надёжное место (это пароль БД и секрет JWT). Если потеряешь — токены сессий сбросятся, но данные останутся.

---

## Шаг 4. Первичная настройка через админку

Открой `http://<ip-сервера>/admin/login`, войди под созданным суперадмином и:

1. **Settings** → загрузи логотип, фавикон, настрой цвета, название института.
2. **Разделы** → создай иерархию разделов (например, «Кафедры» → «Кафедра экономики»).
3. **Книги** → загрузи первые материалы.

Все изменения применяются без перезапуска.

---

## Ежедневные операции

### Просмотр статуса

```bash
cd ~/library-offline-bundle-2026-05-13
docker compose -f docker-compose.prod.yml ps
```

### Логи

```bash
docker compose -f docker-compose.prod.yml logs -f api      # API
docker compose -f docker-compose.prod.yml logs -f frontend # фронт
docker compose -f docker-compose.prod.yml logs -f backup   # ночной бэкап
```

### Перезапуск отдельного сервиса

```bash
docker compose -f docker-compose.prod.yml restart api
```

### Ручной бэкап (вне расписания)

```bash
docker compose -f docker-compose.prod.yml exec backup sh /usr/local/bin/backup.sh
```

Файл появится в `/var/library/backups/daily/library-YYYY-MM-DD.sql.gz`.

---

## Восстановление из бэкапа

```bash
sudo ./scripts/restore.sh
```

Скрипт покажет список бэкапов и попросит подтверждение. Все текущие данные БД будут заменены. Файлы книг (PDF/EPUB) остаются нетронутыми — они лежат в `/var/library/files/`, а в БД хранятся только пути.

Для конкретного файла:

```bash
sudo ./scripts/restore.sh /var/library/backups/weekly/library-2026-19.sql.gz
```

---

## Обновление платформы

> Полный пошаговый сценарий обновления с разбором версий — в [UPDATE.md](UPDATE.md). Здесь — краткая версия.

Когда разработчик пришлёт новый bundle (имя содержит версию, например `library-offline-bundle-1.1.0-2026-08-01.tar.gz`):

```bash
tar xzf library-offline-bundle-1.1.0-2026-08-01.tar.gz
cd library-offline-bundle-1.1.0-2026-08-01

# Копируем существующий .env, чтобы не пересоздавать секреты (пароль БД, JWT)
cp ~/library-offline-bundle-1.0.0-2026-05-13/.env ./

sudo ./scripts/install.sh --update
```

В режиме `--update` установщик только подменяет образы и перезапускает api/frontend/nginx. БД и файлы не трогает. Миграции БД применяются автоматически при старте api. Установщик сам синхронизирует `TAG` в `.env` с версией нового bundle.

**Обязательно проверь, что обновление применилось** — версия в ответе должна совпасть с версией нового bundle:

```bash
curl -s http://127.0.0.1/api/health
# {"status":"ok","version":"1.1.0","db":"connected"}
```

Если версия не сменилась — обновление не доехало (частая причина: запустили `docker compose restart` вместо `install.sh --update`). См. [UPDATE.md](UPDATE.md), раздел «Если фичи не применились».

---

## Удаление

**Только остановка** (данные сохраняются):

```bash
sudo ./scripts/uninstall.sh
```

**Полное удаление** (БД, файлы, бэкапы — всё стирается):

```bash
sudo ./scripts/uninstall.sh --purge
```

---

## Резервное копирование на внешний носитель

Скрипт бэкапит БД внутри сервера. Чтобы защититься от отказа диска, настрой еженедельную копию `/var/library/backups/` на внешний накопитель.

Пример: подключаешь USB-диск в `/mnt/backup-usb`, добавляешь в `/etc/crontab`:

```cron
# Воскресенье, 04:00 — копирование бэкапов на USB
0 4 * * 0 root rsync -av --delete /var/library/backups/ /mnt/backup-usb/library-backups/
```

---

## Безопасность

1. **Firewall**: открой только 80/443 наружу. PostgreSQL и API внутри docker-сети, наружу не торчат.

   ```bash
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw allow ssh
   sudo ufw enable
   ```

2. **HTTPS**: даже в локальной сети. Сгенерируй самоподписанный сертификат:

   ```bash
   sudo openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
     -keyout /var/library/ssl/privkey.pem \
     -out /var/library/ssl/fullchain.pem \
     -subj "/CN=library.institute.local"
   ```

   Затем раскомментируй HTTPS-блок в `nginx/conf.d/library.conf` и перезапусти:

   ```bash
   docker compose -f docker-compose.prod.yml restart nginx
   ```

3. **Пароль БД**: уже сгенерирован случайно при установке. Сменить — поправь `POSTGRES_PASSWORD` в `.env` и:

   ```bash
   docker compose -f docker-compose.prod.yml exec postgres \
     psql -U library -c "ALTER USER library WITH PASSWORD 'NEW_PASSWORD';"
   docker compose -f docker-compose.prod.yml restart api
   ```

---

## Что делать, если PostgreSQL не стартует

1. Проверь логи:

   ```bash
   docker compose -f docker-compose.prod.yml logs postgres
   ```

2. Если повреждение WAL после жёсткого выключения сервера:

   ```bash
   # Восстановление из последнего бэкапа
   sudo systemctl stop docker
   sudo mv /var/library/pg-data /var/library/pg-data.corrupted
   sudo mkdir -p /var/library/pg-data
   sudo chown -R 70:70 /var/library/pg-data
   sudo systemctl start docker
   docker compose -f docker-compose.prod.yml up -d postgres
   # после старта PG:
   sudo ./scripts/restore.sh
   ```

3. Если ничего не помогло — позвони разработчику и пришли архив `pg-data.corrupted/`.

---

## Контакты

Проблемы и обновления — пиши разработчику. Для получения нового bundle понадобится:

- версия платформы (см. `~/library-offline-bundle-*/​docker-compose.prod.yml`, поле `TAG`)
- лог проблемного контейнера
- содержимое `.env` (без `POSTGRES_PASSWORD` и `JWT_SECRET`)
