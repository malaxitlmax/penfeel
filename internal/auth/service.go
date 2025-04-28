package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	pkgauth "github.com/malaxitlmax/penfeel/pkg/auth"
)

// Service интерфейс для сервиса авторизации
type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest) (*TokenResponse, error)
	ValidateToken(ctx context.Context, token string) (*ValidationResponse, error)
}

// AuthService реализация сервиса авторизации
type AuthService struct {
	repo            Repository
	passwordService *pkgauth.PasswordService
	jwtService      *pkgauth.JWTService
}

// NewAuthService создает новый сервис авторизации
func NewAuthService(repo Repository, passwordService *pkgauth.PasswordService, jwtService *pkgauth.JWTService) *AuthService {
	return &AuthService{
		repo:            repo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Хешируем пароль
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаем пользователя
	user := &User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	// Сохраняем пользователя в БД
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login выполняет вход пользователя
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	// Ищем пользователя по email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Проверяем пароль
	if err := s.passwordService.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Генерируем JWT токен
	token, expiresAt, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Генерируем Refresh токен
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Формируем ответ
	response := &TokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         *user,
	}

	return response, nil
}

// ValidateToken проверяет JWT токен
func (s *AuthService) ValidateToken(ctx context.Context, token string) (*ValidationResponse, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return &ValidationResponse{Valid: false}, fmt.Errorf("invalid token: %w", err)
	}

	// Преобразуем строковый ID в UUID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return &ValidationResponse{Valid: false}, fmt.Errorf("invalid user ID: %w", err)
	}

	// Получаем пользователя из БД
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return &ValidationResponse{Valid: false}, fmt.Errorf("user not found: %w", err)
	}

	return &ValidationResponse{
		Valid: true,
		User:  *user,
	}, nil
}
