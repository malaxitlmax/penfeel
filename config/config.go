package config

import (
	"os"
	"strconv"
	"time"
)

// Config структура для хранения конфигурации приложения
type Config struct {
	Database  DatabaseConfig
	JWT       JWTConfig
	Server    ServerConfig
	Migration MigrationConfig
}

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig конфигурация для JWT токенов
type JWTConfig struct {
	Secret          string
	ExpirationHours int
	RefreshSecret   string
	RefreshExpHours int
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	Port         int
	GRPCPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// MigrationConfig конфигурация миграций
type MigrationConfig struct {
	Path             string
	Enabled          bool
	LockTimeout      int
	StatementTimeout int
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "penfeel"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
			RefreshSecret:   getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
			RefreshExpHours: getEnvAsInt("JWT_REFRESH_EXPIRATION_HOURS", 168), // 7 days
		},
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			GRPCPort:     getEnvAsInt("GRPC_PORT", 9090),
			ReadTimeout:  time.Duration(getEnvAsInt("SERVER_READ_TIMEOUT", 10)) * time.Second,
			WriteTimeout: time.Duration(getEnvAsInt("SERVER_WRITE_TIMEOUT", 10)) * time.Second,
			IdleTimeout:  time.Duration(getEnvAsInt("SERVER_IDLE_TIMEOUT", 60)) * time.Second,
		},
		Migration: MigrationConfig{
			Path:             getEnv("MIGRATION_PATH", "./migrations"),
			Enabled:          getEnvAsBool("MIGRATION_ENABLED", true),
			LockTimeout:      getEnvAsInt("MIGRATION_LOCK_TIMEOUT", 5000),
			StatementTimeout: getEnvAsInt("MIGRATION_STATEMENT_TIMEOUT", 60000),
		},
	}
}

// Helper функции для работы с переменными окружения
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
