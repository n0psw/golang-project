# Создание администратора

## Быстрый старт

### Windows (PowerShell)
```powershell
.\scripts\create_admin.ps1
```

### Linux/Mac
```bash
make create-admin
```

или

```bash
./scripts/create_admin.sh
```

## Прямой запуск

```bash
go run cmd/admin/create_admin.go -email admin@example.com -username admin -password yourpassword
```

## Через переменные окружения

```bash
export ADMIN_EMAIL=admin@example.com
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD=yourpassword
go run cmd/admin/create_admin.go
```

## Требования

- База данных должна быть запущена и доступна
- Миграции должны быть применены
- Переменные окружения из `.env` должны быть настроены

## Пример использования

После создания администратора вы можете войти в систему:

```bash
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "yourpassword"
}
```

Ответ будет содержать JWT токен, который можно использовать для доступа к админ-панели.

