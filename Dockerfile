FROM golang:latest as builder
LABEL authors="https://t.me/Zatrasz"

# Создание рабочий директории
RUN mkdir -p /app

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта внутрь контейнера
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o ads ./cmd/main.go

# Второй этап: создание production образ
FROM ubuntu AS chemistry

WORKDIR /app

# Обновляем список пакетов
RUN apt-get update && apt-get install -y nginx

# Копируем исполняемый файл из этапа сборки
COPY --from=builder /app/ads ./

# Копируем Swagger UI статические файлы
COPY ./docs /app/swagger-ui

# Устанавливаем сервер для обслуживания Swagger
COPY ./nginx.conf /etc/nginx/sites-available/default

# Копируем остальное
COPY ./ ./

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]

CMD ["./ads"]