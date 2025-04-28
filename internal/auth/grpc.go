package auth

import (
	"context"

	pb "github.com/malaxitlmax/penfeel/api/proto"
)

// GRPCServer реализация gRPC сервера для авторизации
type GRPCServer struct {
	pb.UnimplementedAuthServiceServer
	service Service
}

// NewGRPCServer создает новый gRPC сервер для авторизации
func NewGRPCServer(service Service) *GRPCServer {
	return &GRPCServer{
		service: service,
	}
}

// Register обрабатывает запрос на регистрацию
func (s *GRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Преобразуем запрос в доменную модель
	domainReq := RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// Вызываем сервис для регистрации
	user, err := s.service.Register(ctx, domainReq)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Формируем ответ
	return &pb.RegisterResponse{
		Success: true,
		UserId:  user.ID.String(),
	}, nil
}

// Login обрабатывает запрос на вход
func (s *GRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Преобразуем запрос в доменную модель
	domainReq := LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Вызываем сервис для входа
	response, err := s.service.Login(ctx, domainReq)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Формируем ответ
	return &pb.LoginResponse{
		Success:      true,
		Token:        response.Token,
		RefreshToken: response.RefreshToken,
		User: &pb.UserInfo{
			Id:       response.User.ID.String(),
			Username: response.User.Username,
			Email:    response.User.Email,
		},
	}, nil
}

// ValidateToken обрабатывает запрос на проверку токена
func (s *GRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// Вызываем сервис для проверки токена
	response, err := s.service.ValidateToken(ctx, req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	// Если токен не валиден, возвращаем ошибку
	if !response.Valid {
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: "invalid token",
		}, nil
	}

	// Формируем ответ
	return &pb.ValidateTokenResponse{
		Valid: true,
		User: &pb.UserInfo{
			Id:       response.User.ID.String(),
			Username: response.User.Username,
			Email:    response.User.Email,
		},
	}, nil
}
