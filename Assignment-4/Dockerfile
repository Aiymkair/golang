# ---- Stage 1: build ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем модули для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Статическая компиляция
RUN go build -o myapp ./cmd/api/main.go

# Stage 2: final image
FROM alpine:latest

# CA-сертификаты для HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/myapp .

COPY --from=builder /app/database ./database

EXPOSE 8080

CMD ["./myapp"]