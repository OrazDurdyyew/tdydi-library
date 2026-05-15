DC := docker compose -f deploy/docker-compose.dev.yml

.PHONY: dev up down restart logs ps admin psql reset graph help

dev: up logs

up:
	$(DC) up -d

down:
	$(DC) down

restart:
	$(DC) restart

logs:
	$(DC) logs -f api frontend

ps:
	$(DC) ps

admin:
	$(DC) exec api go run ./cmd/createadmin -username=admin -password=admin12345 -role=superadmin

psql:
	$(DC) exec postgres psql -U library -d library

reset:
	$(DC) down -v

graph:
	~/.venv/bin/graphify update .

help:
	@echo "make dev      — поднять всё и показать логи (Ctrl+C только остановит логи, контейнеры останутся)"
	@echo "make up       — поднять контейнеры в фоне"
	@echo "make down     — остановить"
	@echo "make restart  — перезапустить"
	@echo "make logs     — логи api и frontend"
	@echo "make ps       — статус контейнеров"
	@echo "make admin    — создать суперадмина admin/admin12345"
	@echo "make psql     — psql в БД"
	@echo "make reset    — снести контейнеры + volumes (потеря данных)"
	@echo "make graph    — обновить graphify"
