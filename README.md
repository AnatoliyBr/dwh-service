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
* [Добавление нового сервиса](#добавление-нового-сервиса)
* [Просмотр сервиса](#просмотр-сервиса)

### Добавление нового сервиса
Добавление нового сервиса:

```bash
curl --location --request POST http://localhost:8080/services \
--data-raw '{
    "slug":"TODO_APP",
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

## Решения
В ходе разработки были сомнения по тем или иным вопросам, которые были решены следующим образом:
1. Как организовать хранение произвольных метрик, набор которых динамически меняется?
* **Первый вариант**. При регистрации события записывать каждую метрику в отдельную строку (то есть одна строка = одна метрика). Схема простая, расширяемая, нормализованная, позволяет использовать реляционную БД типа PostgreSQL и **хорошо справляется с добавлением новых метрик**, но **требует преобразования значений метрик разных типов в один (например, в строку)**.
* **Второй вариант**. Каждой метрике сопоставлять столбец и использовать колоночную БД типа Cassandra. Схема **позволяет эффективнее считывать информацию по ключам**, но **требует изменения схемы БД (добавления новой колонки) при добавлении каждой новой метрики в БД**.

> В данном проекте использовался первый вариант.