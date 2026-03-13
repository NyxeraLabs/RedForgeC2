# Makefile for RedForgeC2 local workflows

COMPOSE := docker compose -f docker/docker-compose.yml

.PHONY: up down ps logs build db-reset clean
.PHONY: teamserver agent ui tui e2e

up:
	$(COMPOSE) up --build

down:
	$(COMPOSE) down

ps:
	$(COMPOSE) ps

logs:
	$(COMPOSE) logs -f --tail=100

build:
	$(COMPOSE) build --no-cache

db-reset:
	$(COMPOSE) down -v

clean:
	docker system prune -af

teamserver:
	cd teamserver && go run ./cmd/teamserver

agent:
	cd agent && cargo run

ui:
	cd ui && npm install && npm run dev

tui:
	cd agent && cargo run --bin tui

e2e:
	./tests/e2e/run.sh

