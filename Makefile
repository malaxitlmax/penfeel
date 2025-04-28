.PHONY: proto build run-auth run-api-gateway migrate docker-build docker-up docker-down

# Генерация proto файлов
proto:
	./scripts/generate_proto.sh

# Сборка всех сервисов
build: proto
	go build -o bin/auth-service cmd/auth/main.go
	go build -o bin/api-gateway cmd/api-gateway/main.go

# Запуск сервиса авторизации
run-auth: build
	./bin/auth-service

# Запуск API Gateway
run-api-gateway: build
	./bin/api-gateway

# Запуск миграций
migrate:
	./scripts/run_migrations.sh

# Сборка Docker образов
docker-build:
	docker compose build

# Запуск Docker контейнеров
docker-up:
	docker compose up -d

# Остановка Docker контейнеров
docker-down:
	docker compose down

# Установка зависимостей для разработки
install-dev-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest 