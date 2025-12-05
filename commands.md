# Команды для запуска проекта (без Docker)

## Требования
- Go 1.24+
- Node.js 20+
- PostgreSQL 15+ (или Docker для PostgreSQL)
- golang-migrate CLI

---

## Backend (Go API)

### 1. Запустить PostgreSQL

**Вариант A: Через Docker (рекомендуется)**
```bash
docker run -d --name postgres-movie -e POSTGRES_PASSWORD=password -e POSTGRES_DB=movie_review_db -p 5432:5432 postgres:15-alpine
```

**Вариант B: Использовать существующий контейнер**
```bash
docker start movie-review-api-postgres-1
```

**Вариант C: Локальный PostgreSQL**
- Убедитесь, что PostgreSQL запущен
- Создайте базу: `CREATE DATABASE movie_review_db;`

### 2. Настроить переменные окружения
```bash
cd movie-review-api
cp env.example .env
```

Отредактируйте `.env` если нужно (по умолчанию работает с настройками выше).

### 3. Применить миграции
```bash
cd movie-review-api
make migrate-up
```

Или вручную:
```bash
migrate -path migrations -database "postgres://postgres:password@localhost:5432/movie_review_db?sslmode=disable" up
```

### 4. Запустить API сервер
```bash
cd movie-review-api
make run
```

Или:
```bash
cd movie-review-api
go run cmd/api/main.go
```

✅ API будет доступен на: `http://localhost:8080`

---

## Frontend (React)

### 1. Установить зависимости (если еще не установлены)
```bash
cd frontend
npm install
```

### 2. Запустить dev сервер
```bash
cd frontend
npm run dev
```

✅ Frontend будет доступен на: `http://localhost:5173`

Vite автоматически проксирует запросы `/api/*` на `http://localhost:8080/api/*`

---

## Быстрый старт (оба сервиса)

### Терминал 1 - Backend:
```bash
cd movie-review-api
go run cmd/api/main.go
```

### Терминал 2 - Frontend:
```bash
cd frontend
npm run dev
```

Откройте браузер: `http://localhost:5173`

---

## Создание админа

```bash
cd movie-review-api
go run cmd/admin/create_admin.go -email admin@test.com -username admin -password 123123
```

---

## Остановка

- **Backend**: `Ctrl+C` в терминале
- **Frontend**: `Ctrl+C` в терминале  
- **PostgreSQL**: `docker stop postgres-movie` (если через Docker)

---

## Проверка работы

1. Backend: откройте `http://localhost:8080/api/v1/genres` - должны увидеть список жанров
2. Frontend: откройте `http://localhost:5173` - должна открыться главная страница
