
# Web Makefile - convenient targets for local development, build and docker

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || date -u +"%Y%m%d%H%M%S")
IMAGE_NAME := terarush/web
COMPOSE := docker-compose -f docker-compose.yml

.PHONY: help build clean dev preview docker-build docker-up docker-down version install full-install

help:
	@printf "Available targets:\n"
	@printf "  build          Install deps and build production assets\n"
	@printf "  clean          Run project clean script (and remove build cache)\n"
	@printf "  dev            Run vite dev server\n"
	@printf "  preview        Preview the production build (vite preview)\n"
	@printf "  docker-build   Build docker image $(IMAGE_NAME):$(VERSION)\n"
	@printf "  docker-up      docker-compose up -d --build\n"
	@printf "  docker-down    docker-compose down\n"
	@printf "  version        Print derived version string\n"
	@printf "  install        Copy built assets to /var/www/terarush and record version (sudo)\n"

build:
	@echo "Building web assets"
	@npm ci
	@npm run build

clean:
	@echo "Cleaning web assets"
	@npm run clean || true
	@rm -rf node_modules/.tmp

dev:
	@echo "Starting dev server"
	@npm run dev

preview:
	@echo "Previewing production build"
	@npm run preview

docker-build:
	@echo "Building docker image $(IMAGE_NAME):$(VERSION)"
	@docker build -t $(IMAGE_NAME):$(VERSION) .

docker-up:
	@echo "Bringing up docker-compose services"
	@$(COMPOSE) up -d --build

docker-down:
	@echo "Stopping docker-compose services"
	@$(COMPOSE) down

version:
	@printf "%s\n" "$(VERSION)"

install: build
	@echo "Installing built assets to /var/www/terarush (requires sudo)"
	@sudo mkdir -p /var/www/terarush
	@sudo rm -rf /var/www/terarush/*
	@sudo cp -r dist/* /var/www/terarush/
	@printf "%s\n" "$(VERSION)" | sudo tee /etc/terarush_web_version >/dev/null

vercel:
	@echo "Deploying to Vercel"
	@vercel deploy --prod

vercel-dev:
	@echo "Deploying to Vercel dev environment"
	@vercel deploy --prod --environment dev

full-install: install
