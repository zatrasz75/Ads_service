---
title: README
tags: []
---
# **Тестовое задание: REST API для хранения и подачи объявлений**

## **Описание**

Реализовать сервис для хранения и подачи объявлений с использованием языка программирования Go. Сервис должен предоставлять API, работающее поверх HTTP в формате JSON.

## **Требования**

- Язык программирования: Go
- Финальная версия должна быть выложена на GitHub.
- Простая инструкция для запуска, предпочтительно с возможностью запустить через `docker-compose up`.
- Реализация трех методов: получение списка объявлений, получение одного объявления, создание объявления.

## **Детали**

### **Метод получения списка объявлений**

- Пагинация: на одной странице должно присутствовать 10 объявлений.
- Сортировки: по цене (возрастание/убывание) и по дате создания (возрастание/убывание).
- Поля в ответе: название объявления, цена.

### **Метод получения конкретного объявления**

- Обязательные поля в ответе: название объявления, цена.
- Опциональные поля (можно запросить, передав параметр `fields`): описание.

### **Метод создания объявления**

- Принимает все вышеперечисленные поля: название, описание, цена.
- Возвращает ID созданного объявления и код результата (ошибка или успех).

## **Усложнения**

- Юнит тесты: постарайтесь достичь покрытия в 70% и больше.
- Контейнеризация: есть возможность поднять проект с помощью команды `docker-compose up`.
- Документация: Swagger.

## **Инструкция по запуску**

1. Клонируйте репозиторий:

```
git clone https://github.com/zatrasz75/Ads_service.git
```

1. Перейдите в директорию проекта:

```
cd Ads_service
```

1. Запустите проект с помощью Docker Compose:

```
docker-compose up
```

Запуск сервера на <http://localhost:3131>

Документация Swagger API: <http://localhost:3131/swagger/index.html>

### endpoints:

/posts \[GET\] Получение конкретного объявления по ID

/posts \[POST\] Создание нового объявления

/posts/list \[GET\] Получение списка объявлений

1. Запустите проект на компьютере, предварительно установив Golang и MongoDB настроив MONGO_CONN_STR в  .env:

установить зависимости

```
go mod download
go mod tidy
```

запуск

```
go run cmd/main.go
```

Запуск сервера на <http://localhost:3232>

Документация Swagger API: <http://localhost:3232/swagger/index.html>

### endpoints:

/posts \[GET\] Получение конкретного объявления по ID

/posts \[POST\] Создание нового объявления

/posts/list \[GET\] Получение списка объявлений

## **Вопросы и принятые решения**

- **Какие поля будут в объявлении?**\: Название, описание, цена.
- **Как будет реализована пагинация?**\: Используя параметры запроса для указания номера страницы и 10 объявлений на странице.
- Если параметр не задан - по умолчанию 1я страница
- **Как будет реализована сортировка?**\: Используя параметры запроса для указания полей фильтрации направления сортировки.
- **Как будет реализована фильтрация?**\: Используя параметры запроса для указания критериев фильтрации цена и дата создания.