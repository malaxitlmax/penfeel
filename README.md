# PenFeel - Сервис для писателей с коллаборативным редактированием

PenFeel - это веб-приложение, позволяющее писателям совместно работать над текстами в реальном времени.

## Технологический стек

- Backend: Golang
- База данных: PostgreSQL
- Коммуникация между сервисами: gRPC
- Frontend: React (планируется)
- Деплой: Docker, docker-compose

## Структура проекта

```
/
├── cmd/                      # Точки входа для всех сервисов
│   ├── auth/                 # Сервис авторизации
│   └── ...                   # Другие сервисы (будут добавлены)
├── pkg/                      # Общие пакеты для переиспользования
├── internal/                 # Код, специфичный для сервисов
├── api/                      # API-контракты и документация
├── migrations/               # Миграции БД
├── scripts/                  # Скрипты для CI/CD и деплоя
├── config/                   # Конфигурационные файлы
```

## Начало работы

### Предварительные требования

- Go 1.23+
- Docker и docker-compose
- PostgreSQL (для локальной разработки)
- protoc (для генерации Proto файлов)

### Установка зависимостей для разработки

```bash
make install-dev-deps
```

### Генерация Proto-файлов

```bash
make proto
```

### Запуск с помощью Docker

```bash
# Сборка образов
make docker-build

# Запуск контейнеров
make docker-up

# Запуск миграций
make migrate

# Остановка контейнеров
make docker-down
```

### Локальный запуск для разработки

```bash
# Запуск сервиса авторизации
make run-auth
```

## Сервисы

### Auth Service (Сервис авторизации)

- Регистрация и вход пользователей
- Управление JWT-токенами
- Проверка прав доступа
