#!/bin/bash
set -e

# Проверяем, установлен ли migrate
if ! command -v migrate &> /dev/null; then
    echo "Error: migrate is not installed"
    echo "Install with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# Собираем строку подключения к базе данных
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-penfeel}
DB_SSLMODE=${DB_SSLMODE:-disable}

DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

echo "Running migrations on ${DB_URL}..."

# Запускаем миграции
migrate -source file://migrations -database "${DB_URL}" up

echo "Migrations completed successfully!" 