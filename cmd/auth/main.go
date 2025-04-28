package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/malaxitlmax/penfeel/config"
	"github.com/malaxitlmax/penfeel/internal/auth"
	pkgauth "github.com/malaxitlmax/penfeel/pkg/auth"
	"github.com/malaxitlmax/penfeel/pkg/database"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Подключаемся к базе данных
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Создаем репозиторий
	repo := auth.NewPostgresRepository(db)

	// Создаем сервисы
	passwordService := pkgauth.NewPasswordService(bcrypt.DefaultCost)
	jwtService := pkgauth.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.ExpirationHours,
		cfg.JWT.RefreshSecret,
		cfg.JWT.RefreshExpHours,
	)

	// Создаем сервис авторизации
	authService := auth.NewAuthService(repo, passwordService, jwtService)

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем сервис авторизации
	authGRPCServer := auth.NewGRPCServer(authService)
	// TODO: Раскомментировать после генерации proto
	// pb.RegisterAuthServiceServer(grpcServer, authGRPCServer)

	// Включаем reflection для отладки с помощью grpcurl
	reflection.Register(grpcServer)

	// Запускаем gRPC сервер
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Auth service starting on port %d", cfg.Server.GRPCPort)

	// Обрабатываем сигналы для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Ждем сигнала для завершения
	<-stop

	log.Println("Shutting down Auth service")
	grpcServer.GracefulStop()
}
