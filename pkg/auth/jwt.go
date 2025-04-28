package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService сервис для работы с JWT токенами
type JWTService struct {
	secretKey     string
	expirationHrs int
	refreshKey    string
	refreshExpHrs int
}

// Claims структура данных для JWT токена
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// NewJWTService создает новый JWT сервис
func NewJWTService(secretKey string, expirationHrs int, refreshKey string, refreshExpHrs int) *JWTService {
	return &JWTService{
		secretKey:     secretKey,
		expirationHrs: expirationHrs,
		refreshKey:    refreshKey,
		refreshExpHrs: refreshExpHrs,
	}
}

// GenerateToken генерирует JWT токен для пользователя
func (s *JWTService) GenerateToken(userID uuid.UUID) (string, time.Time, error) {
	expirationTime := time.Now().Add(time.Duration(s.expirationHrs) * time.Hour)

	claims := &Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expirationTime, nil
}

// GenerateRefreshToken генерирует refresh токен
func (s *JWTService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.refreshExpHrs) * time.Hour)

	claims := &Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.refreshKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken проверяет валидность токена
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// ValidateRefreshToken проверяет валидность refresh токена
func (s *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.refreshKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	return claims, nil
}
