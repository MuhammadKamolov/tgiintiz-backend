# Этап 1: Сборка бинарного файла
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Сначала копируем зависимости, чтобы Docker кэшировал этот слой
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и собираем статический бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/api ./cmd/api

# Этап 2: Минимальный production-образ
FROM alpine:3.21

WORKDIR /app

# Сертификаты для HTTPS и часовые пояса
RUN apk --no-cache add ca-certificates tzdata

# Создаём непривилегированного пользователя для безопасности
RUN adduser -D -g '' appuser \
    && mkdir -p /app/uploads \
    && chown -R appuser:appuser /app

COPY --from=builder /app/api /app/api

USER appuser

EXPOSE 8080

CMD ["/app/api"]