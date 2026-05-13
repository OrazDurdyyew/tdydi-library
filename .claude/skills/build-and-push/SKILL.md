---
name: build-and-push
description: Собрать Docker-образы api и frontend и запушить в приватный registry. Триггер — когда пользователь говорит "запушь образы", "собери Docker", "выпусти релиз", "обнови инстанс ВУЗа". Тегирует версией, собирает multi-stage, пушит.
---

# build-and-push

Собирает и публикует Docker-образы в приватный registry для деплоя у ВУЗов.

## Когда использовать

- "запушь новую версию"
- "собери образы"
- "выпусти релиз"
- "обнови инстанс"

## Что делает

1. Спросить у пользователя версию тега (например, `v1.2.0` или `latest`). По умолчанию — `latest` + `git rev-parse --short HEAD`.
2. Проверить что `git status` чистый. Если нет — спросить подтверждение.
3. Собрать два образа:
   ```sh
   docker build -t $REGISTRY/library-api:$TAG -t $REGISTRY/library-api:latest ./api
   docker build -t $REGISTRY/library-frontend:$TAG -t $REGISTRY/library-frontend:latest ./frontend
   ```
4. Запушить оба:
   ```sh
   docker push $REGISTRY/library-api:$TAG
   docker push $REGISTRY/library-api:latest
   docker push $REGISTRY/library-frontend:$TAG
   docker push $REGISTRY/library-frontend:latest
   ```
5. Создать git-тег: `git tag $TAG && git push origin $TAG`

## Переменные окружения

- `$REGISTRY` — например, `ghcr.io/username` или `harbor.example.com/library`
- Логин в registry должен быть выполнен заранее: `docker login $REGISTRY`

## Соглашения

- Версии — semver: `vMAJOR.MINOR.PATCH`
- `latest` всегда указывает на самый последний релиз
- Не пушить `latest` без явного подтверждения пользователя

## После push

- Напомнить пользователю команду обновления на сервере ВУЗа:
  ```sh
  ssh vuz-server "cd /opt/library && docker compose pull && docker compose up -d"
  ```
- Опционально: пингнуть health endpoint после деплоя

## Перед первым использованием

- Убедиться что Dockerfile в `/api/` и `/frontend/` есть и multi-stage
- Образ `api`: alpine + Go-бинарь + `pdftoppm` (нужен для генерации обложек)
- Образ `frontend`: build stage с node, runtime — node:alpine с собранным `.output/`
