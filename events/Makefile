APP_NAME := terarush-events
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || date +%Y%m%d%H%M%S)
COMPOSE  := docker compose --env-file ../.env

.PHONY: build run up down restart logs ps deploy help

build:
	@echo "Building $(APP_NAME) binary..."
	@go build -o bin/$(APP_NAME) .
	@echo "Done → bin/$(APP_NAME)"

up:
	@echo "Starting $(APP_NAME) (version: $(VERSION))..."
	@$(COMPOSE) up -d --force-recreate

down:
	@$(COMPOSE) down

restart:
	@$(COMPOSE) restart events

logs:
	@$(COMPOSE) logs -f events

ps:
	@$(COMPOSE) ps

deploy: down
	@echo "Building $(APP_NAME) (version: $(VERSION))..."
	@$(COMPOSE) build --no-cache
	@$(COMPOSE) up -d --force-recreate
	@echo "$(APP_NAME) deployed (version: $(VERSION))"

help:
	@echo ""
	@echo "  $(APP_NAME) — Makefile targets"
	@echo ""
	@echo "  make build    - build Go binary locally"
	@echo "  make up       - build image & start container"
	@echo "  make down     - stop & remove container"
	@echo "  make restart  - restart running container"
	@echo "  make logs     - tail container logs"
	@echo "  make ps       - show container status"
	@echo "  make deploy   - down then up (full redeploy)"
	@echo ""
