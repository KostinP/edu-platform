# Stage 1: Build
FROM golang:1.24.3 AS builder

WORKDIR /app
COPY . .
RUN go build -o app ./cmd/main.go

# Stage 2: Runtime
FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/app .

# Пример: копируем миграции, если есть
# COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./app"]
