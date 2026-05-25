build:
	go build -o bin/app main.go
	echo "Build completed. Executable created at bin/app"

run:
	docker compose up --build -d
	npm run build --prefix web && npm run start --prefix web

install:
	make build
	mv bin/app /usr/local/bin/app
	cp .env /usr/local/bin
	cp config.service /etc/systemd/system/app.service
	systemctl daemon-reload
	systemctl enable app.service
	systemctl start app.service

install-docker:
	docker compose up --build -d