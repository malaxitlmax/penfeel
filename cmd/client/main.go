package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/malaxitlmax/penfeel/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Определяем адрес gRPC сервера с возможностью переопределения через переменную окружения
	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = "localhost:9090"
	}

	// Подключаемся к gRPC серверу
	log.Printf("Connecting to %s", addr)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Создаем клиент для сервиса авторизации
	client := pb.NewAuthServiceClient(conn)

	// Устанавливаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Выбираем операцию в зависимости от аргументов командной строки
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s [register|login|validate]", os.Args[0])
	}

	switch os.Args[1] {
	case "register":
		if len(os.Args) != 5 {
			log.Fatalf("Usage: %s register <username> <email> <password>", os.Args[0])
		}
		username, email, password := os.Args[2], os.Args[3], os.Args[4]
		register(ctx, client, username, email, password)
	case "login":
		if len(os.Args) != 4 {
			log.Fatalf("Usage: %s login <email> <password>", os.Args[0])
		}
		email, password := os.Args[2], os.Args[3]
		login(ctx, client, email, password)
	case "validate":
		if len(os.Args) != 3 {
			log.Fatalf("Usage: %s validate <token>", os.Args[0])
		}
		token := os.Args[2]
		validateToken(ctx, client, token)
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}

func register(ctx context.Context, client pb.AuthServiceClient, username, email, password string) {
	log.Printf("Registering user: %s (%s)", username, email)

	resp, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Fatalf("Failed to register user: %v", err)
	}

	if resp.Success {
		log.Printf("User registered successfully. User ID: %s", resp.UserId)
	} else {
		log.Printf("Failed to register user: %s", resp.Error)
	}
}

func login(ctx context.Context, client pb.AuthServiceClient, email, password string) {
	log.Printf("Logging in user: %s", email)

	resp, err := client.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	if resp.Success {
		log.Printf("Login successful")
		log.Printf("Token: %s", resp.Token)
		log.Printf("Refresh Token: %s", resp.RefreshToken)
		log.Printf("User: ID=%s, Username=%s, Email=%s",
			resp.User.Id, resp.User.Username, resp.User.Email)
	} else {
		log.Printf("Login failed: %s", resp.Error)
	}
}

func validateToken(ctx context.Context, client pb.AuthServiceClient, token string) {
	log.Printf("Validating token")

	resp, err := client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}

	if resp.Valid {
		log.Printf("Token is valid")
		log.Printf("User: ID=%s, Username=%s, Email=%s",
			resp.User.Id, resp.User.Username, resp.User.Email)
	} else {
		log.Printf("Token is invalid: %s", resp.Error)
	}
}
