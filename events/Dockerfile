FROM golang:1.24-alpine

RUN apk add --no-cache \
    ca-certificates \
    make \
    git \
    bash \
    docker-cli \
    docker-cli-compose \
    openssh-client

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o events-service .

CMD ["./events-service"]
