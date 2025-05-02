package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL драйвер
	"github.com/malaxitlmax/penfeel/config"
	"github.com/malaxitlmax/penfeel/internal/database/migration"
)

// NewPostgresDB создает новое подключение к базе данных PostgreSQL
func NewPostgresDB(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	return db, nil
}

// NewPostgresDBWithMigrations creates a database connection and runs migrations
func NewPostgresDBWithMigrations(cfg config.DatabaseConfig, migrationPath string) (*sqlx.DB, error) {
	// Create database connection string
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	// Run migrations
	if migrationPath != "" {
		log.Println("Running database migrations...")
		migrationCfg := migration.Config{
			MigrationsPath:   migrationPath,
			DatabaseURL:      dsn,
			LockTimeout:      5000,  // 5 seconds
			StatementTimeout: 60000, // 60 seconds
		}

		if err := migration.RunMigrations(migrationCfg); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	// Connect to database
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	return db, nil
}
