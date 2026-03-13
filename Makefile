# Makefile for RedForgeC2 development / docker workflows

DOCKER_COMPOSE := docker-compose -f docker/docker-compose.yml

.PHONY: all up down build restart logs ps clean db-reset agent teamserver ui

all: up

# Bring up the full lab (postgres + teamserver + ui)
up:
	$(DOCKER_COMPOSE) up --build

# Stop the lab (leaves volumes intact)
down:
	$(DOCKER_COMPOSE) down

# Stop and remove volumes (reset state)
db-reset:
	$(DOCKER_COMPOSE) down -v

# Rebuild images without bringing them up
build:
	$(DOCKER_COMPOSE) build --no-cache

# View logs
logs:
	$(DOCKER_COMPOSE) logs -f --tail=100

# List running containers
ps:
	$(DOCKER_COMPOSE) ps

# Clean up unused docker images/volumes
clean:
	docker system prune -af

# Restart the lab (down + up in one command)
restart: down up

# Run just the agent locally (not inside docker)
agent:
	cd agent && cargo run

# Run just the teamserver locally (not inside docker)
teamserver:
	cd teamserver && go run ./cmd/teamserver

# Run just the UI locally (not inside docker)
ui:
	cd ui && npm install && npm run dev
