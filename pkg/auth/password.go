package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService сервис для работы с паролями
type PasswordService struct {
	cost int
}

// NewPasswordService создает новый сервис для паролей
func NewPasswordService(cost int) *PasswordService {
	// Если передан некорректный cost, используем значение по умолчанию
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}

	return &PasswordService{
		cost: cost,
	}
}

// HashPassword хеширует пароль
func (s *PasswordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword проверяет правильность пароля
func (s *PasswordService) CheckPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	return nil
}
