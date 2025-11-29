FROM golang:1.24.2-alpine3.20 AS builder

RUN apk update && apk add --no-cache \
    git \
    openssh \
    tzdata \
    build-base \
    ca-certificates \
    curl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/grpc-server ./cmd/grpc/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/rest-server ./cmd/rest/main.go

FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata ca-certificates curl && \
    mkdir -p /app

WORKDIR /app

EXPOSE 50052 3000

COPY --from=builder /app/bin/grpc-server /app/grpc-server
COPY --from=builder /app/bin/rest-server /app/rest-server
COPY --from=builder /app/storage /app/storage

RUN mkdir -p /app/storage/product

# Create empty .env file so godotenv.Load() doesn't fail
# Actual environment variables are injected by docker-compose at runtime
RUN touch /app/.env

CMD ["/app/grpc-server"]