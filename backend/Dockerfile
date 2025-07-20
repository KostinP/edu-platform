# Stage 1: Build
FROM golang:1.24.3 AS builder

WORKDIR /app

# Копируем модули для кеширования слоёв
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

# Stage 2: Runtime
FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
