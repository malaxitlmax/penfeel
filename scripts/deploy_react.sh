#!/bin/bash

# Скрипт для копирования собранного React-приложения в нужную директорию
# Выполнять из корневой директории проекта

# Проверяем, существует ли директория для статических файлов
mkdir -p client/dist

# Собираем React-приложение
cd client && npm run build

echo "React app built successfully!"

# Запуск API Gateway
echo "Starting API Gateway..."
cd ..
go run cmd/api-gateway/main.go 