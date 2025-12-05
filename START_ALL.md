# Запуск всего проекта одной командой

## 🚀 Быстрый старт

### Запустить все сервисы:

```bash
cd movie-review-api
docker compose up -d --build
```

Это запустит:
- ✅ PostgreSQL (база данных)
- ✅ Миграции (автоматически)
- ✅ API сервер (Go)
- ✅ Фронтенд (React + Nginx)

### Доступ к приложению:

- **Фронтенд:** http://localhost:3000
- **API:** http://localhost:8080/api/v1

### Остановить все:

```bash
cd movie-review-api
docker compose down
```

### Перезапустить:

```bash
cd movie-review-api
docker compose restart
```

### Просмотр логов:

```bash
# Все сервисы
docker compose logs -f

# Только API
docker compose logs -f api

# Только фронтенд
docker compose logs -f frontend
```

## 🔧 Локальная разработка (без Docker для фронтенда)

Если нужен hot reload для фронтенда:

```bash
# 1. Запустить бэкенд в Docker
cd movie-review-api
docker compose up -d postgres migrate api

# 2. Запустить фронтенд локально (в другом терминале)
cd frontend
npm run dev
```

Доступ:
- Фронтенд: http://localhost:5173 (с hot reload)
- API: http://localhost:8080

## 📋 Проверка работы

### Проверить API:
```bash
curl http://localhost:8080/api/v1/genres
```

### Проверить статус контейнеров:
```bash
docker compose ps
```

### Проверить логи:
```bash
docker compose logs api
docker compose logs frontend
```

## 🐛 Решение проблем

### Порт занят:
- Измените порты в `docker-compose.yml`
- Или остановите другие сервисы на этих портах

### Ошибки сборки:
```bash
docker compose build --no-cache
docker compose up -d
```

### Очистить все и начать заново:
```bash
docker compose down -v
docker compose up -d --build
```

## 📝 Примечания

- Первый запуск может занять время (скачивание образов, сборка)
- Миграции применяются автоматически при первом запуске
- Данные сохраняются в Docker volume `postgres_data`
- Фронтенд в Docker работает в production mode (статический сервер `serve`)
- Для разработки с hot reload используйте локальный запуск: `cd frontend && npm run dev`

