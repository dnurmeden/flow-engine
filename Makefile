ENV_FILE ?= .env
include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))

DB_DSN := $(DATABASE_URL)

.PHONY: up down logs psql migrate-up migrate-down migrate-status new-migration wait-db lint test help

help:
	@echo "Targets: up, down, logs, psql, migrate-up, migrate-down, migrate-status, new-migration, wait-db, lint, test"

up:
	docker compose -f deployments/docker-compose.yml --env-file $(ENV_FILE) up -d

down:
	docker compose -f deployments/docker-compose.yml --env-file $(ENV_FILE) down --remove-orphans

nuke: ## остановить и удалить контейнер, сеть и volume
	docker compose -f deployments/docker-compose.yml --env-file $(ENV_FILE) down --volumes --remove-orphans || true
	-docker rm -f flow-postgres 2>/dev/null || true
	-docker network prune -f

reset-db: nuke up migrate-up

logs:
	docker compose -f deployments/docker-compose.yml logs -f

# Ждем здоровья Postgres перед миграциями
wait-db:
	@echo "Waiting for Postgres..."; \
	for i in $$(seq 1 30); do \
	  docker inspect -f '{{.State.Health.Status}}' flow-postgres 2>/dev/null | grep -q healthy && echo "OK" && exit 0; \
	  sleep 1; \
	done; \
	echo "Postgres is not healthy" && exit 1

psql:
	docker exec -it flow-postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

migrate-up: wait-db
	goose -dir ./migrations postgres "$(DB_DSN)" up

migrate-down:
	goose -dir ./migrations postgres "$(DB_DSN)" down

migrate-status:
	goose -dir ./migrations postgres "$(DB_DSN)" status

new-migration:
	@read -p "Migration name: " name; \
	goose -dir ./migrations create $$name sql

lint:
	golangci-lint run

test:
	go test ./... -race -count=1
