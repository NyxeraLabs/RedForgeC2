# Makefile for RedForgeC2 local workflows

COMPOSE := docker compose --project-directory $(CURDIR) --env-file $(CURDIR)/.env -f $(CURDIR)/docker/docker-compose.yml

.PHONY: up down ps logs build db-reset clean
.PHONY: env-init env-check
.PHONY: teamserver agent ui tui e2e

env-init:
	@if [ ! -f $(CURDIR)/.env ]; then \
		cp $(CURDIR)/.env.example $(CURDIR)/.env; \
		echo "Created .env from .env.example. Edit .env and set required values before running 'make up'."; \
	fi

env-check:
	@set -eu; \
	if [ ! -f $(CURDIR)/.env ]; then \
		echo "Missing .env. Run: cp .env.example .env (then edit it)"; \
		exit 1; \
	fi; \
	required="REDFORGE_DB_PASS REDFORGE_JWT_SECRET REDFORGE_ADMIN_PASS"; \
	for k in $$required; do \
		v="$${!k-}"; \
		if [ -z "$$v" ]; then \
			v="$$(awk -F= -v key="$$k" 'BEGIN{found=0} $$1==key{found=1; sub(/^[^=]*=/,""); print; exit} END{if(!found)exit 2}' $(CURDIR)/.env 2>/dev/null || true)"; \
		fi; \
		if [ -z "$$v" ] || [ "$$v" = "change_me" ] || [ "$$v" = "change_me_to_a_long_random_value" ] || [ "$$v" = "change_me_to_a_strong_password" ]; then \
			echo "$$k is missing/placeholder. Set it in .env (or export it) and retry."; \
			exit 1; \
		fi; \
	done

up: env-init env-check
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
