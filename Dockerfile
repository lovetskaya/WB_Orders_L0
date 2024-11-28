# Builder stage
FROM golang:1.22.3-alpine AS builder

WORKDIR /usr/local/src

# Копируем для управления зависимостями
COPY ["go.mod", "go.sum", "./"]

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . ./

# Сборка Go-приложения
RUN go build -o ./bin/app cmd/app/main.go

# с Go для тестов
FROM golang:1.22.3-alpine AS runner

# Установка необходимых зависимостей
RUN apk add --no-cache ca-certificates

# Копируем скомпилированное приложение из builder stage
COPY --from=builder /usr/local/src/bin/app /

# Копируем весь исходный код для запуска тестов
COPY . /usr/local/src/

WORKDIR /usr/local/src/

EXPOSE 8080

CMD ["/app"]