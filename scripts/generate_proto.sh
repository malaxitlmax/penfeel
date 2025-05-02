#!/bin/sh
set -e

# Генерируем код из proto файлов
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  api/proto/auth.proto \
  api/proto/document.proto

echo "Proto files generated successfully" 