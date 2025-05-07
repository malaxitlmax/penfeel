.PHONY: proto run-auth run-api migrate build up down migration-new migration-up migration-down migration-status migration-plan dump-schema

# Генерация proto файлов
proto:
	./scripts/generate_proto.sh

# Запуск сервиса авторизации
run-auth: build
	./bin/auth-service

# Запуск API
run-api: build
	./bin/api

# Запуск миграций (legacy)
migrate:
	./scripts/run_migrations.sh

# Migration commands (new)
migration-new:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $${name}

# Database URL construction
define get_db_url
	$(eval DB_HOST ?= localhost)
	$(eval DB_PORT ?= 5432)
	$(eval DB_USER ?= postgres)
	$(eval DB_PASSWORD ?= postgres)
	$(eval DB_NAME ?= penfeel)
	$(eval DB_SSLMODE ?= disable)
	$(eval DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE))
endef

migration-up:
	$(call get_db_url)
	@echo "Running migrations on $(DB_URL)..."
	migrate -path migrations -database "$(DB_URL)" up

migration-down:
	$(call get_db_url)
	@read -p "How many migrations to roll back? " num; \
	echo "Rolling back $$num migrations on $(DB_URL)..."; \
	migrate -path migrations -database "$(DB_URL)" down $$num

migration-status:
	$(call get_db_url)
	@echo "Current migration version on $(DB_URL):"
	migrate -path migrations -database "$(DB_URL)" version

migration-plan:
	$(call get_db_url)
	@echo "The following migrations will be applied to $(DB_URL):"
	@migrate -path migrations -database "$(DB_URL)" up -dry-run

# Schema dump
dump-schema:
	./scripts/dump_schema.sh

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
	docker exec -it penfeel-api sh -c "cd client && npm install"

npm-build:
	docker exec -it penfeel-api sh -c "cd client && npm run build"

build-base:
	docker build -t penfeel-base -f Dockerfile.base .