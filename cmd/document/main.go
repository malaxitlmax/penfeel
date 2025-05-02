package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	pb "github.com/malaxitlmax/penfeel/api/proto"
	"github.com/malaxitlmax/penfeel/config"
	"github.com/malaxitlmax/penfeel/internal/document"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()
	isDev := os.Getenv("ENV") == "dev"
	_ = isDev

	// Настраиваем логирование
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Подключаемся к базе данных
	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Создаем репозиторий
	repo := document.NewPostgresRepository(db)

	// Создаем сервис
	service := document.NewDocumentService(repo)

	// Создаем gRPC сервер
	grpcServer := document.NewGRPCServer(service)

	// Запускаем gRPC сервер
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterDocumentServiceServer(server, grpcServer)

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Document service starting on port %d", cfg.Server.GRPCPort)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Настройка graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Document service...")
	server.GracefulStop()
	log.Println("Document service stopped")
}
