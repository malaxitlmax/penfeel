package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/malaxitlmax/penfeel/config"
	authGateway "github.com/malaxitlmax/penfeel/internal/api-gateway/auth"
	documentGateway "github.com/malaxitlmax/penfeel/internal/api-gateway/document"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/malaxitlmax/penfeel/api/proto"
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

	// Получаем хост auth-service из переменной окружения или используем localhost по умолчанию
	authServiceHost := os.Getenv("AUTH_SERVICE_HOST")
	if authServiceHost == "" {
		authServiceHost = "localhost"
	}

	// Получаем хост document-service из переменной окружения или используем localhost по умолчанию
	documentServiceHost := os.Getenv("DOCUMENT_SERVICE_HOST")
	if documentServiceHost == "" {
		documentServiceHost = "localhost"
	}

	// Получаем порт document-service из переменной окружения или используем стандартный
	documentServicePort := cfg.Server.GRPCPort
	if portStr := os.Getenv("DOCUMENT_SERVICE_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			documentServicePort = port
		}
	}

	// Устанавливаем соединение с auth service через gRPC
	authConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", authServiceHost, cfg.Server.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	// Устанавливаем соединение с document service через gRPC
	documentConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", documentServiceHost, documentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to document service: %v", err)
	}
	defer documentConn.Close()

	// Создаем клиенты для сервисов
	authClient := pb.NewAuthServiceClient(authConn)
	documentClient := pb.NewDocumentServiceClient(documentConn)

	// Создаем роутер gin
	router := gin.Default()

	if isDev {
		// Configure CORS middleware
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowAllOrigins = true // For development; restrict in production
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
		corsConfig.ExposeHeaders = []string{"Content-Length"}
		corsConfig.AllowCredentials = true
		router.Use(cors.New(corsConfig))
	}

	// Регистрируем маршруты для аутентификации
	authHandler := authGateway.NewHandler(authClient)
	authMiddleware := authGateway.AuthMiddleware(authClient)

	// Регистрируем маршруты для документов
	documentHandler := documentGateway.NewHandler(documentClient)

	// Путь к собранному React-приложению
	staticPath := "./client/dist"

	// Публичные маршруты
	authRoutes := router.Group("/api/v1/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/validate", authHandler.ValidateToken)
	}

	// Защищенные маршруты (пример)
	protectedRoutes := router.Group("/api/v1")
	protectedRoutes.Use(authMiddleware)
	{
		// Пример защищенного маршрута
		protectedRoutes.GET("documents", documentHandler.GetDocuments)
		protectedRoutes.GET("documents/:id", documentHandler.GetDocument)
		protectedRoutes.POST("documents", documentHandler.CreateDocument)
		protectedRoutes.PUT("documents/:id", documentHandler.UpdateDocument)
		protectedRoutes.DELETE("documents/:id", documentHandler.DeleteDocument)
	}

	// TODO: включать на проде
	// Обслуживание статических файлов (должно быть последним)
	// router.Static("/", staticPath)

	// Настройка и запуск HTTP сервера
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Запуск сервера в горутине
	go func() {
		log.Printf("API Gateway starting on port %d", cfg.Server.Port)
		log.Printf("Serving React app from %s", staticPath)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Настройка graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("API Gateway stopped")
}
