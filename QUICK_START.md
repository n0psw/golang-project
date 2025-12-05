# Быстрый запуск проекта

## 1. Запуск бэкенда

### Вариант A: Через Docker (рекомендуется)

1. **Запустите Docker Desktop**

2. **Запустите базу данных, миграции и API:**
```bash
cd movie-review-api
docker compose up -d
```

Это автоматически:
- Запустит PostgreSQL
- Применит миграции базы данных
- Запустит API сервер

3. **Проверьте, что все запущено:**
```bash
docker compose ps
```

4. **Проверьте API:**
```bash
curl http://localhost:8080/api/v1/genres
```

### Вариант B: Локально (если PostgreSQL установлен)

1. **Убедитесь, что PostgreSQL запущен**

2. **Создайте базу данных:**
```sql
CREATE DATABASE movie_review_db;
```

3. **Примените миграции:**
```bash
cd movie-review-api
make migrate-up
```

Если `make` не работает:
```bash
migrate -path migrations -database "postgres://postgres:password@localhost:5432/movie_review_db?sslmode=disable" up
```

4. **Запустите API:**
```bash
go run cmd/api/main.go
```

## 2. Запуск фронтенда

```bash
cd frontend
npm run dev
```

Фронтенд будет доступен на: http://localhost:5173

## 3. Проверка работы

- Бэкенд: http://localhost:8080/api/v1
- Фронтенд: http://localhost:5173

## Решение проблем

### Ошибка подключения к базе данных

**Если используете Docker:**
- Убедитесь, что Docker Desktop запущен
- Проверьте: `docker-compose ps`

**Если используете локальный PostgreSQL:**
- Убедитесь, что PostgreSQL запущен
- Проверьте настройки в `.env` файле
- Убедитесь, что база данных создана

### Ошибка миграций

Установите golang-migrate:
- Windows: `choco install golang-migrate`
- Или скачайте: https://github.com/golang-migrate/migrate/releases

### Порт 8080 занят

Измените `PORT` в `movie-review-api/.env` и обновите `API_BASE_URL` в `frontend/src/api/client.ts`

