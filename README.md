# Задача

Добавить возможность задавать настройки периодичности задачи.

# Основные изменения
- У сущности задачи Task добавлено поле настройки периодичности Recurrence, у которого есть поле типа задачи PeriodicType:
  - NotPeriodic (непериодическая)
  - Daily (каждый n-й день)
  - Monthly (заданное число каждого месяца)
  - ExactDates (конкретные даты)
  - EvenOddDates (четные / нечетные дни)

- Настройки для каждого типа задаются соответствующими полями в Recurrence:
  - EveryNDays (раз в сколько дней повторяется задачи, если тип Daily)
  - MonthlyDate (число месяца от 1 до 30, если тип Monthly)
  - Dates (список дат, если тип ExactDates)
  - EvenOddDates (определяет повтор по четным / нечетным дням, если тип EvenOddDates)

- Обновлены методы работы с БД postgres
- Добавлена валидация входных данных для периодических задач
- Обновлены методы обработчиков запросов
- Обновлена документация Swagger
- Добавлены unit тесты для use case слоя и для HTTP обработчиков

## Task Service

Сервис для управления задачами с HTTP API на Go.

## Требования

- Go `1.23+`
- Docker и Docker Compose

## Быстрый запуск через Docker Compose

```bash
docker compose up --build
```

После запуска сервис будет доступен по адресу `http://localhost:8080`.

Если `postgres` уже запускался ранее со старой схемой, пересоздай volume:

```bash
docker compose down -v
docker compose up --build
```

Причина в том, что SQL-файл из `migrations/0001_create_tasks.up.sql` монтируется в `docker-entrypoint-initdb.d` и применяется только при инициализации пустого data volume.

## Swagger

Swagger UI:

```text
http://localhost:8080/swagger/
```

OpenAPI JSON:

```text
http://localhost:8080/swagger/openapi.json
```

## API

Базовый префикс API:

```text
/api/v1
```

Основные маршруты:

- `POST /api/v1/tasks`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/{id}`
- `PUT /api/v1/tasks/{id}`
- `DELETE /api/v1/tasks/{id}`
