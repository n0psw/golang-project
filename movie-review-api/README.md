# Movie Review API

REST API для системы отзывов о фильмах, построенная на Go. Позволяет пользователям просматривать фильмы, оставлять отзывы и управлять контентом (для администраторов).

## Технологии

- **Go 1.24+** - основной язык программирования
- **PostgreSQL** - база данных
- **Gin** - веб-фреймворк
- **JWT** - аутентификация
- **Docker** - контейнеризация
- **golang-migrate** - миграции базы данных

## Требования

- Go 1.24 или выше
- PostgreSQL 15+
- Docker и Docker Compose (опционально)
- golang-migrate CLI

## Установка и запуск

### 1. Клонирование репозитория

```bash
git clone <repository-url>
cd movie-review-api
```

### 2. Настройка окружения

Создайте файл `.env` на основе `env.example`:

```bash
cp env.example .env
```

Отредактируйте `.env` с вашими настройками:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=movie_review_db
DB_SSLMODE=disable

JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRES_IN=24h

PORT=8080
HOST=0.0.0.0

LOG_LEVEL=info
```

### 3. Установка зависимостей

```bash
go mod download
```

### 4. Настройка базы данных

Убедитесь, что PostgreSQL запущен и создана база данных:

```sql
CREATE DATABASE movie_review_db;
```

### 5. Применение миграций

```bash
make migrate-up
```

Или вручную:

```bash
migrate -path migrations -database "postgres://postgres:password@localhost:5432/movie_review_db?sslmode=disable" up
```

### 6. Создание администратора

Для создания администратора используйте один из способов:

**Способ 1: Через Makefile (Linux/Mac)**
```bash
make create-admin
```

**Способ 2: Через скрипт (Linux/Mac)**
```bash
chmod +x scripts/create_admin.sh
./scripts/create_admin.sh
```

**Способ 3: Через PowerShell (Windows)**
```powershell
.\scripts\create_admin.ps1
```

**Способ 4: Напрямую через Go**
```bash
go run cmd/admin/create_admin.go -email admin@example.com -username admin -password yourpassword
```

**Способ 5: Через переменные окружения**
```bash
export ADMIN_EMAIL=admin@example.com
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD=yourpassword
go run cmd/admin/create_admin.go
```

### 7. Запуск приложения

#### Локальный запуск

```bash
make run
```

Или:

```bash
go run cmd/api/main.go
```

#### Docker Compose

```bash
docker-compose up --build
```

Приложение будет доступно по адресу `http://localhost:8080`

## API Endpoints

Базовый URL: `http://localhost:8080/api/v1`

### Аутентификация

#### Регистрация
```
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "testuser",
  "password": "password123"
}
```

#### Вход
```
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

Ответ:
```json
{
  "user": {
    "id": "...",
    "email": "user@example.com",
    "username": "testuser",
    "role": "user"
  },
  "token": "jwt-token-here"
}
```

#### Получить профиль
```
GET /users/me
Authorization: Bearer <token>
```

### Фильмы

#### Получить все фильмы
```
GET /movies?page=1&limit=10&genre=Action&year=2020&min_rating=7.5&search=matrix
```

Параметры запроса:
- `page` - номер страницы (по умолчанию 1)
- `limit` - количество на странице (по умолчанию 10)
- `genre` - фильтр по жанру
- `year` - фильтр по году
- `min_rating` - минимальный рейтинг
- `search` - поиск по названию/описанию

#### Получить фильм по ID
```
GET /movies/:id
```

#### Создать фильм (только admin)
```
POST /movies
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "title": "The Matrix",
  "description": "A computer hacker learns about the true nature of reality",
  "release_year": 1999,
  "director": "Lana Wachowski",
  "duration_minutes": 136,
  "genre_ids": ["genre-uuid-1", "genre-uuid-2"]
}
```

#### Обновить фильм (только admin)
```
PUT /movies/:id
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "title": "Updated Title",
  "description": "Updated description",
  "release_year": 2000,
  "director": "Director Name",
  "duration_minutes": 150,
  "genre_ids": ["genre-uuid-1"]
}
```

#### Удалить фильм (только admin)
```
DELETE /movies/:id
Authorization: Bearer <admin-token>
```

### Отзывы

#### Создать отзыв
```
POST /movies/:id/reviews
Authorization: Bearer <token>
Content-Type: application/json

{
  "rating": 8,
  "title": "Great movie!",
  "content": "I really enjoyed this film..."
}
```

#### Получить отзывы фильма
```
GET /movies/:id/reviews?page=1&limit=10
```

#### Получить отзыв по ID
```
GET /reviews/:id
Authorization: Bearer <token>
```

#### Получить мои отзывы
```
GET /users/me/reviews?page=1&limit=10
Authorization: Bearer <token>
```

#### Обновить отзыв
```
PUT /reviews/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "rating": 9,
  "title": "Updated title",
  "content": "Updated content"
}
```

#### Удалить отзыв
```
DELETE /reviews/:id
Authorization: Bearer <token>
```

### Жанры

#### Получить все жанры
```
GET /genres
```

#### Получить жанр по ID
```
GET /genres/:id
```

#### Создать жанр (только admin)
```
POST /genres
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Thriller"
}
```

#### Обновить жанр (только admin)
```
PUT /genres/:id
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Updated Genre Name"
}
```

#### Удалить жанр (только admin)
```
DELETE /genres/:id
Authorization: Bearer <admin-token>
```

## Аутентификация

API использует JWT (JSON Web Tokens) для аутентификации. После успешного входа или регистрации вы получите токен, который нужно передавать в заголовке `Authorization`:

```
Authorization: Bearer <your-jwt-token>
```

Токен действителен 24 часа (настраивается через `JWT_EXPIRES_IN`).

## Роли пользователей

- **user** - обычный пользователь, может создавать и управлять своими отзывами
- **admin** - администратор, имеет доступ ко всем операциям, включая управление фильмами и жанрами

## Тестирование

Запуск тестов:

```bash
make test
```

Или:

```bash
go test ./...
```

## Структура проекта

```
movie-review-api/
├── cmd/
│   └── api/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── database/
│   │   └── database.go         # Подключение к БД
│   ├── handler/
│   │   ├── auth_handler.go     # Обработчики аутентификации
│   │   ├── movie_handler.go    # Обработчики фильмов
│   │   ├── review_handler.go   # Обработчики отзывов
│   │   └── genre_handler.go    # Обработчики жанров
│   ├── middleware/
│   │   └── middleware.go       # Middleware (CORS, Auth, Admin)
│   ├── models/
│   │   └── models.go           # Модели данных
│   ├── repository/
│   │   ├── user_repository.go  # Репозиторий пользователей
│   │   ├── movie_repository.go # Репозиторий фильмов
│   │   ├── review_repository.go# Репозиторий отзывов
│   │   └── genre_repository.go # Репозиторий жанров
│   └── service/
│       ├── auth_service.go     # Сервис аутентификации
│       ├── movie_service.go    # Сервис фильмов
│       ├── review_service.go    # Сервис отзывов
│       ├── genre_service.go    # Сервис жанров
│       └── rating_worker.go    # Background worker для обновления рейтингов
├── migrations/
│   ├── 000001_create_tables.up.sql
│   ├── 000001_create_tables.down.sql
│   ├── 000002_seed_data.up.sql
│   └── 000002_seed_data.down.sql
├── pkg/
│   └── jwt/
│       └── jwt.go              # JWT утилиты
├── tests/
│   ├── auth_service_test.go    # Тесты сервиса аутентификации
│   └── integration_test.go     # Интеграционные тесты
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

## Особенности

- **Concurrency**: Background worker (RatingWorker) обновляет средние рейтинги фильмов асинхронно при создании/обновлении/удалении отзывов
- **Graceful Shutdown**: Приложение корректно завершает работу при получении сигналов SIGINT/SIGTERM
- **Context Propagation**: Все операции используют context для отмены и таймаутов
- **Structured Logging**: Логирование в формате JSON через slog
- **Dependency Injection**: Использование интерфейсов для тестируемости

## Makefile команды

- `make migrate-up` - применить миграции
- `make migrate-down` - откатить миграции
- `make migrate-create name=migration_name` - создать новую миграцию
- `make test` - запустить тесты
- `make build` - собрать приложение
- `make run` - запустить приложение
- `make docker-build` - собрать Docker образ
- `make docker-run` - запустить Docker контейнер

## Лицензия

MIT

