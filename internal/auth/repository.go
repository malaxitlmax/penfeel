package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repository интерфейс для работы с хранилищем пользователей
type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
}

// PostgresRepository реализация репозитория для PostgreSQL
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository создает новый репозиторий для PostgreSQL
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// CreateUser создает нового пользователя в базе данных
func (r *PostgresRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	err := r.db.QueryRowxContext(
		ctx,
		query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
	).StructScan(user)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail получает пользователя по email
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetUserByID получает пользователя по ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

// ErrUserNotFound ошибка, когда пользователь не найден
var ErrUserNotFound = errors.New("user not found")
