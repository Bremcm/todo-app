# Todo App — REST API на Go

REST API для управления списками задач: регистрация и аутентификация пользователей по JWT,
CRUD над списками (TodoList) и задачами внутри них (TodoItem). Учебный pet-проект, написанный
для практики backend-разработки на Go.

## Стек

- **Go** — основной язык
- **Gin** — HTTP-фреймворк / роутинг
- **PostgreSQL** — база данных
- **sqlx** — работа с БД
- **golang-migrate** — миграции схемы
- **Viper** — конфигурация (`configs/config.yml`)
- **godotenv** — секреты из `.env`
- **JWT** (`dgrijalva/jwt-go`) — аутентификация
- **bcrypt** — хеширование паролей
- **logrus** — логирование
- **Docker / docker-compose** — запуск

## Архитектура

Приложение разбито на три слоя с однонаправленным потоком зависимостей:

```
handler  →  service  →  repository  →  PostgreSQL
(HTTP)      (логика)     (SQL)
```

- **handler** — парсинг запроса, вызов сервиса, формирование ответа
- **service** — бизнес-логика и валидация (например, проверка, что список принадлежит пользователю)
- **repository** — SQL-запросы к Postgres

Маршрут запроса проходит сверху вниз, ответ — обратно. JWT проверяется middleware до попадания
в защищённые хендлеры.

## Структура проекта

```
.
├── cmd/main.go                 # точка входа, graceful shutdown
├── configs/config.yml          # несекретная конфигурация
├── schema/                     # SQL-миграции (up / down)
├── pkg/handler/                # HTTP-слой
├── pkg/handler/service/        # бизнес-логика
├── pkg/handler/service/repository/  # работа с БД
├── server.go                   # обёртка http.Server
├── todo.go / user.go           # доменные структуры
├── Dockerfile
├── docker-compose.yml
└── .env                        # секреты (не коммитится)
```

## Конфигурация

Несекретные параметры — в `configs/config.yml`:

```yaml
port: "8000"

db:
  host: "localhost"      # для docker-compose поменяйте на "db"
  port: "5436"
  username: "postgres"
  dbname: "postgres"
  sslmode: "disable"
```

Секреты — в `.env` (скопируйте из `.env.example`):

```
DB_PASSWORD=qwerty
SIGNING_KEY=change_me_to_a_long_random_secret
```

## Запуск

### Вариант 1 — через docker-compose (рекомендуется)

Поднимает Postgres, накатывает миграции и запускает приложение одной командой.

```bash
cp .env.example .env        # и впишите свои значения
docker-compose up --build
```

> Внутри Docker приложение обращается к БД по имени сервиса `db`.
> Убедитесь, что в `config.yml` для этого режима `host: "db"`.

API будет доступен на `http://localhost:8000`.

### Вариант 2 — локально (go run)

1. Поднимите Postgres в Docker:

```bash
docker run --name todo-db -e POSTGRES_PASSWORD=qwerty -p 5436:5432 -d postgres:15-alpine
```

2. Накатите миграции (нужен установленный `migrate`):

```bash
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up
```

3. Запустите приложение (`host: "localhost"` в конфиге):

```bash
go run cmd/main.go
```

## API

### Аутентификация

| Метод | Путь            | Описание                       | Auth |
|-------|-----------------|--------------------------------|------|
| POST  | `/auth/sign-up` | Регистрация пользователя       | —    |
| POST  | `/auth/sign-in` | Логин, возвращает JWT-токен    | —    |

### Списки (требуют заголовок `Authorization: Bearer <token>`)

| Метод  | Путь              | Описание                |
|--------|-------------------|-------------------------|
| POST   | `/api/lists/`     | Создать список          |
| GET    | `/api/lists/`     | Все списки пользователя |
| GET    | `/api/lists/:id`  | Список по id            |
| PUT    | `/api/lists/:id`  | Обновить список         |
| DELETE | `/api/lists/:id`  | Удалить список          |

### Задачи

| Метод  | Путь                                | Описание             |
|--------|-------------------------------------|----------------------|
| POST   | `/api/lists/:id/items/`             | Создать задачу       |
| GET    | `/api/lists/:id/items/`             | Все задачи списка    |
| GET    | `/api/lists/:id/items/:item_id`     | Задача по id         |
| PUT    | `/api/lists/:id/items/:item_id`     | Обновить задачу      |
| DELETE | `/api/lists/:id/items/:item_id`     | Удалить задачу       |

## Примеры

Регистрация:

```bash
curl -X POST http://localhost:8000/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{"name":"Brem","username":"bremcm","password":"qwerty"}'
```

Логин (возвращает токен):

```bash
curl -X POST http://localhost:8000/auth/sign-in \
  -H "Content-Type: application/json" \
  -d '{"username":"bremcm","password":"qwerty"}'
```

Создание списка (подставьте токен из логина):

```bash
curl -X POST http://localhost:8000/api/lists/ \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"My list","description":"Things to do"}'
```

## Схема БД

Пять таблиц: `users`, `todo_lists`, `todo_items` и две связующие — `users_lists`
(пользователь ↔ список) и `lists_items` (список ↔ задача). Связи — many-to-many через
внешние ключи с `ON DELETE CASCADE`.

## Возможные улучшения

- Реорганизовать пакеты в плоскую структуру (`pkg/handler`, `pkg/service`, `pkg/repository`)
- Перейти с устаревшего `dgrijalva/jwt-go` на `golang-jwt/jwt`
- Возвращать `401` вместо `500` при неверных кредах
- Покрыть слои unit-тестами
