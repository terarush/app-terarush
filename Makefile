APP_NAME := $(shell basename $(CURDIR))
# VERSION: prefer git describe tags, fallback to timestamp
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || date +%Y%m%d%H%M%S)

.PHONY: build run install-local install full-install web-build docker-build docker-up docker-down version help clean

build:
	go build -o bin/app main.go
	@echo "Build completed. Executable created at bin/app"

run:
	docker compose up --build -d
	@# For local dev you can still run the web dev server (not used in production)
	@npm run dev --prefix web

install-local:
    make build
    @sudo mv bin/app /usr/local/bin/$(APP_NAME)

	@printf "[Unit]\n" | sudo tee /etc/systemd/system/$(APP_NAME).service > /dev/null
	@printf "Description=%s Service\n" "$(APP_NAME)" | sudo tee -a /etc/systemd/system/$(APP_NAME).service > /dev/null
	@printf "After=network.target\n\n" | sudo tee -a /etc/systemd/system/$(APP_NAME).service > /dev/null
	@printf "[Service]\n" | sudo tee -a /etc/systemd/system/$(APP_NAME).service > /dev/null
	@printf "ExecStart=/usr/local/bin/%s\n" "$(APP_NAME)" | sudo tee -a /etc/systemd/system/$(APP_NAME).service > /dev/null
	@printf "WorkingDirectory=%s\n" "$(CURDIR)" | sudo tee -a /etc/systemd/system/$(APP_NAME).service > /dev/null
	@printf "Restart=always\n" | sudo tee -a /etc/systemd/system/$(APP_NAME).service > /dev/null
	@printf "RestartSec=5\n\n" | sudo tee -a /etc/systemd/system/$(APP_NAME).service > /dev/null
	@printf "[Install]\nWantedBy=multi-user.target\n" | sudo tee -a /etc/systemd/system/$(APP_NAME).service > /dev/null
    @sudo systemctl daemon-reload
    @sudo systemctl enable $(APP_NAME).service
    @sudo systemctl start $(APP_NAME).service

## Build web static assets and image
web-build:
	@echo "Building web assets (version: $(VERSION))"
	@npm ci --prefix web
	@npm run build --prefix web
	@docker build -t $(APP_NAME)-web:$(VERSION) -f web/Dockerfile web

## Build backend image (optional, if you want a single image)
docker-build:
	@echo "Building backend image (version: $(VERSION))"
	@docker build -t $(APP_NAME):$(VERSION) .

## Start using docker compose (will use local docker-compose.yml)
docker-up:
	@echo "Starting containers via docker compose"
	@docker compose up -d --build

docker-down:
	@docker compose down

## Full install: build web, build images, run compose and persist version
full-install: web-build docker-build docker-up
	@echo "$(VERSION)" | sudo tee /etc/$(APP_NAME)_version > /dev/null
	@echo "Deployed $(APP_NAME) version $(VERSION)"

install: full-install

version:
	@echo $(VERSION)

help:
	@echo "Makefile targets:"
	@echo "  make build            - build Go binary"
	@echo "  make web-build        - build web assets and docker image"
	@echo "  make docker-build     - build backend docker image"
	@echo "  make docker-up        - start containers with docker compose"
	@echo "  make full-install     - build everything and deploy (writes /etc/$(APP_NAME)_version)"

clean:
	rm -rf bin dist
