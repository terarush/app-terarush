build:
	go build -o bin/app main.go
	echo "Build completed. Executable created at bin/app"

run:
	docker compose up --build -d
	npm run build --prefix web && npm run start --prefix web

install:
	make build
	mv bin/app /usr/local/bin/app

install-docker:
	docker compose up --build -d