# Инструкция по запуску бэкенда

## Вариант 1: Запуск через Docker Compose (рекомендуется)

### 1. Убедитесь, что Docker Desktop запущен

### 2. Запустите базу данных и API:

```bash
cd movie-review-api
docker-compose up -d
```

Это запустит:
- PostgreSQL на порту 5432
- API на порту 8080

### 3. Примените миграции:

```bash
docker-compose exec api migrate -path /app/migrations -database "postgres://postgres:password@postgres:5432/movie_review_db?sslmode=disable" up
```

Или если миграции не применяются автоматически, выполните локально:

```bash
migrate -path migrations -database "postgres://postgres:password@localhost:5432/movie_review_db?sslmode=disable" up
```

## Вариант 2: Локальный запуск (без Docker)

### 1. Установите PostgreSQL

Убедитесь, что PostgreSQL установлен и запущен.

### 2. Создайте базу данных:

```sql
CREATE DATABASE movie_review_db;
```

### 3. Настройте .env файл:

Файл `.env` уже создан в `movie-review-api/`. Проверьте настройки:

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

### 4. Установите golang-migrate:

Windows:
```bash
choco install golang-migrate
```

Или скачайте с https://github.com/golang-migrate/migrate/releases

### 5. Примените миграции:

```bash
cd movie-review-api
make migrate-up
```

Или вручную:
```bash
migrate -path migrations -database "postgres://postgres:password@localhost:5432/movie_review_db?sslmode=disable" up
```

### 6. Запустите API:

```bash
cd movie-review-api
make run
```

Или:
```bash
go run cmd/api/main.go
```

## Проверка работы

После запуска API должен быть доступен на:
- http://localhost:8080/api/v1

Проверьте здоровье API:
```bash
curl http://localhost:8080/api/v1/health
```

## Остановка

Если используете Docker:
```bash
docker-compose down
```

## Решение проблем

### Ошибка подключения к базе данных
- Убедитесь, что PostgreSQL запущен
- Проверьте настройки в .env файле
- Проверьте, что база данных создана

### Ошибка миграций
- Убедитесь, что golang-migrate установлен
- Проверьте строку подключения к базе данных

### Порт 8080 занят
- Измените PORT в .env файле
- Обновите API_BASE_URL во frontend/src/api/client.ts

