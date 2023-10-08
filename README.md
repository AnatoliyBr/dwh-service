## Data Warehouse Service
### Описание
REST API приложение для хранения **произвольных метрик**, способное принимать в себя данные, хранить и отдавать пользователю.

### Функционал
Данное серверное (backend) приложение предоставляет **API** с форматом **JSON** для добавления **сервисов**, **метрик** и **событий**, а также получения данных по **ключам событий** за заданный интервал времени.

Под событием понимается получение **списка метрик** от **одного сервиса**.

## Структура REST API

```
POST /services - добавление нового сервиса
GET /services - просмотр отслеживаемых сервисов

POST /metrics - добавление новой метрики
GET /metrics - просмотр используемых метрик

POST /events - добавление нового события
GET /events - получение данных по идентификатору сервиса и метрики за заданный интервал времени
```

## Схема базы данных

<p align="center">
    <img src="/assets/images/db_schema.png" width="800">
</p>

## Архитектура

<p align="center">
    <img src="/assets/images/architecture.png" width="800">
</p>

## Запуск и отладка
Все команды, используемые в процессе разработки и тестирования, фиксировались в `Makefile`.

Чтобы поднять проект, необходимо выполнить **две задачи** из `Makefile`:

```bash
make compose-build
make compose-up
```

## Примеры запросов
* [Добавление сервиса](#добавление-сервиса)
* [Просмотр сервиса](#просмотр-сервиса)
* [Добавление метрики](#добавление-метрики)
    * [Типа INT](#типа-int)
    * [Типа FLOAT](#типа-float)
    * [Типа DURATION](#типа-duration)
    * [Типа TIMESTAMP_WITH_TIMEZONE](#типа-timestamp_with_timezone)
    * [Типа BOOL](#типа-bool)
    * [Типа STRING](#типа-string)
* [Просмотр метрики](#просмотр-метрики)
* [Добавление события](#добавление-события)
* [Получение данных](#получение-данных)

### Добавление сервиса
Добавление нового сервиса:

```bash
curl --location --request POST http://localhost:8080/services \
--data-raw '{
    "slug": "TODO_APP",
    "details":"REST API application for managing task lists (todo lists)"
}'
```

Пример ответа:

```bash
{
    "service_id": 1,   
    "slug": "TODO_APP",
    "details": "REST API application for managing task lists (todo lists)"
}
```

### Просмотр сервиса
Просмотр сервиса по идентификатору:

```bash
curl --location --request GET http://localhost:8080/services \
--data-raw '{
    "service_id": 1
}'
```

Пример ответа:

```bash
{
    "service_id": 1,
    "slug": "TODO_APP",
    "details": "REST API application for managing task lists (todo lists)"
}
```

### Добавление метрики
> Сервер поддерживает 6 типов: "INT", "FLOAT", "DURATION", "TIMESTAMP_WITH_TIMEZONE", "BOOL", "STRING".

#### Типа INT
Добавление новой метрики типа "INT":

```bash
curl --location --request POST http://localhost:8080/metrics \
--data-raw '{
    "slug": "INT_METRIC",
    "metric_type": "INT",
    "details":"Calculated in integers"
}'
```

Пример ответа:

```bash
{
    "metric_id": 1,
    "slug": "INT_METRIC",
    "metric_type": "INT",
    "details": "Calculated in integers"
}
```

#### Типа FLOAT
Добавление новой метрики типа "FLOAT":

```bash
curl --location --request POST http://localhost:8080/metrics \
--data-raw '{
    "slug": "FLOAT_METRIC",
    "metric_type": "FLOAT",
    "details":"Calculated in floating point numbers"
}'
```

Пример ответа:

```bash
{
    "metric_id": 2,
    "slug": "FLOAT_METRIC",
    "metric_type": "FLOAT",
    "details": "Calculated in floating point numbers"
}
```

#### Типа DURATION
Добавление новой метрики типа "DURATION":

```bash
curl --location --request POST http://localhost:8080/metrics \
--data-raw '{
    "slug": "DURATION_METRIC",
    "metric_type": "DURATION",
    "details":"Calculated by duration"
}'
```

Пример ответа:

```bash
{
    "metric_id": 3,
    "slug": "DURATION_METRIC",
    "metric_type": "DURATION",
    "details": "Calculated by duration"
}
```

#### Типа TIMESTAMP_WITH_TIMEZONE
Добавление новой метрики типа "TIMESTAMP_WITH_TIMEZONE":

```bash
curl --location --request POST http://localhost:8080/metrics \
--data-raw '{
    "slug": "TIMESTAMP_WITH_TIMEZONE_METRIC",
    "metric_type": "TIMESTAMP_WITH_TIMEZONE",
    "details":"Calculated by timestamps with timezone"
}'
```

Пример ответа:

```bash
{
    "metric_id": 4,
    "slug": "TIMESTAMP_WITH_TIMEZONE_METRIC",
    "metric_type": "TIMESTAMP_WITH_TIMEZONE",
    "details": "Calculated by timestamps with timezone"
}
```

#### Типа BOOL
Добавление новой метрики типа "BOOL":

```bash
curl --location --request POST http://localhost:8080/metrics \
--data-raw '{
    "slug": "BOOL_METRIC",
    "metric_type": "BOOL",
    "details":"Calculated by logical type"
}'
```

Пример ответа:

```bash
{
    "metric_id": 5,
    "slug": "BOOL_METRIC",
    "metric_type": "BOOL",
    "details": "Calculated by logical type"
}
```

#### Типа STRING
Добавление новой метрики типа "STRING":

```bash
curl --location --request POST http://localhost:8080/metrics \
--data-raw '{
    "slug": "STRING_METRIC",
    "metric_type": "STRING",
    "details":"Contains a message"
}'
```

Пример ответа:

```bash
{
    "metric_id": 6,
    "slug": "STRING_METRIC",
    "metric_type": "STRING",
    "details": "Contains a message"
}
```

### Просмотр метрики
Просмотр метрики по идентификатору:

```bash
curl --location --request GET http://localhost:8080/metrics \
--data-raw '{
    "metric_id": 1
}'
```

Пример ответа:

```bash
{
    "metric_id": 1,
    "slug": "INT_METRIC",
    "metric_type": "INT",
    "details": "Calculated in integers"
}
```

### Добавление события
Добавление нового события со списком метрик:

```bash
curl --location --request POST http://localhost:8080/events \
--data-raw '{
    "service_id": 1,
    "metrics": [
        {
            "metric_id": 1,
            "metric_value": 25
        },
        {
            "metric_id": 2,
            "metric_value": 36.6
        },
        {
            "metric_id": 3,
            "metric_value": "1h13m0.5s"
        },
        {
            "metric_id": 4,
            "metric_value": "2023-10-08T19:14:05Z"
        },
        {
            "metric_id": 5,
            "metric_value": true
        },
        {
            "metric_id": 6,
            "metric_value": "Suspicious activity"
        }
    ]
}'
```

Пример ответа:

```bash
{
    "event": {
        "event_id": 1,
        "time_stamp": "2023-10-08T20:44:59+03:00",
        "service_id": 1
    },
    "metrics": [
        {
            "metric_id": 1,
            "metric_value": 25
        },
        {
            "metric_id": 2,
            "metric_value": 36.6
        },
        {
            "metric_id": 3,
            "metric_value": "1h13m0.5s"
        },
        {
            "metric_id": 4,
            "metric_value": "2023-10-08T19:14:05Z"
        },
        {
            "metric_id": 5,
            "metric_value": true
        },
        {
            "metric_id": 6,
            "metric_value": "Suspicious activity"
        }
    ]
}
```

### Получение данных
Получение данных по идентификатору сервиса и метрики за заданный интервал времени:

```bash
curl --location --request GET http://localhost:8080/events \
--data-raw '{
    "service_id": 1,
    "period": [
        "2023-10-06T10:00:00+03:00",
        "2023-10-09T10:00:00+03:00"
        ],
    "metric_id": 3
}'
```

Пример ответа:

```bash
{
    "request": {
        "service_id": 1,
        "period": [
            "2023-10-06T10:00:00+03:00",
            "2023-10-09T10:00:00+03:00"
        ],
        "metric_id": 3
    },
    "report": [
        {
            "time_stamp": "2023-10-08T20:41:37Z",
            "value": "3m10s"
        },
        {
            "time_stamp": "2023-10-08T20:42:02Z",
            "value": "1h3m9.9s"
        },
        {
            "time_stamp": "2023-10-08T20:42:19Z",
            "value": "33m35s"
        },
        {
            "time_stamp": "2023-10-08T20:42:36Z",
            "value": "13m2.1s"
        }
    ]
}
```

## Решения
В ходе разработки были сомнения по тем или иным вопросам, которые были решены следующим образом:
1. Как организовать хранение произвольных метрик, набор которых динамически меняется?
* **Первый вариант**. При регистрации события записывать каждую метрику в отдельную строку (то есть одна строка = одна метрика). Схема простая, расширяемая, нормализованная, позволяет использовать реляционную БД типа PostgreSQL и **хорошо справляется с добавлением новых метрик**, но **требует преобразования значений метрик разных типов в один (например, в строку)**.
* **Второй вариант**. Каждой метрике сопоставлять столбец и использовать колоночную БД типа Cassandra. Схема **позволяет эффективнее считывать информацию по ключам**, но **требует изменения схемы БД (добавления новой колонки) при добавлении каждой новой метрики в БД**.

> В данном проекте использовался первый вариант.