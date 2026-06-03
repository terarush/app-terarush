FROM golang:1.24-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o events-service .

FROM alpine:3.20

RUN apk add --no-cache \
    ca-certificates \
    make \
    git \
    bash \
    docker-cli \
    docker-cli-compose

WORKDIR /app

COPY --from=builder /build/events-service .

CMD ["./events-service"]
