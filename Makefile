.PHONY: proto run-auth run-api-gateway migrate build up down

# Генерация proto файлов
proto:
	./scripts/generate_proto.sh

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
build:
	docker compose build

# Запуск Docker контейнеров
up:
	docker compose up -d

# Остановка Docker контейнеров
down:
	docker compose down

# Просмотр логов
logs:
	docker compose logs -f

# Установка зависимостей для разработки
install-dev-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest 

npm-install:
	docker exec -it penfeel-api-gateway sh -c "cd client && npm install"

npm-build:
	docker exec -it penfeel-api-gateway sh -c "cd client && npm run build"

