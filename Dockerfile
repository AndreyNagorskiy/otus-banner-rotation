FROM golang:1.23 as build

WORKDIR /app

# Копируем файлы модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и билдим приложение
COPY . .

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
    -o banner-rotation \
    -ldflags "$LDFLAGS" \
    ./cmd/server/*

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/banner-rotation /app/
COPY ./internal/storage/migrations /app/internal/storage/migrations

CMD ["./banner-rotation"]