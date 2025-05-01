package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/malaxitlmax/penfeel/config"
	apigateway "github.com/malaxitlmax/penfeel/internal/api-gateway/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/malaxitlmax/penfeel/api/proto"
)

func main() {
	godotenv.Load()
	isDev := os.Getenv("ENV") == "dev"
	_ = isDev

	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Получаем хост auth-service из переменной окружения или используем localhost по умолчанию
	authServiceHost := os.Getenv("AUTH_SERVICE_HOST")
	if authServiceHost == "" {
		authServiceHost = "localhost"
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

	// Создаем клиент для auth service
	authClient := pb.NewAuthServiceClient(authConn)

	// Создаем роутер gin
	router := gin.Default()

	// Регистрируем маршруты для аутентификации
	authHandler := apigateway.NewHandler(authClient)
	authMiddleware := apigateway.AuthMiddleware(authClient)

	// Путь к собранному React-приложению
	staticPath := "./client/dist"

	// Публичные маршруты
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/validate", authHandler.ValidateToken)
	}

	// Защищенные маршруты (пример)
	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(authMiddleware)
	{
		// Пример защищенного маршрута
		protectedRoutes.GET("/user-info", func(c *gin.Context) {
			userId := c.GetString("user_id")
			username := c.GetString("username")
			email := c.GetString("email")

			c.JSON(http.StatusOK, gin.H{
				"user_id":  userId,
				"username": username,
				"email":    email,
			})
		})
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
