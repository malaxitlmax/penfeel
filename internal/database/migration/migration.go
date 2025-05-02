package migration

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Config holds configuration for migrations
type Config struct {
	// MigrationsPath is the path to migration files
	MigrationsPath string
	// DatabaseURL is the connection string for the database
	DatabaseURL string
	// LockTimeout is the timeout for acquiring a lock in milliseconds
	LockTimeout int
	// StatementTimeout is the timeout for SQL statements in milliseconds
	StatementTimeout int
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() Config {
	return Config{
		MigrationsPath:   "./migrations",
		LockTimeout:      5000,  // 5 seconds
		StatementTimeout: 60000, // 60 seconds
	}
}

// RunMigrations applies all pending migrations
func RunMigrations(config Config) error {
	// Validate config
	if config.DatabaseURL == "" {
		return errors.New("database URL is required")
	}

	if config.MigrationsPath == "" {
		config.MigrationsPath = "./migrations"
	}

	// Check if migrations directory exists
	absPath, err := filepath.Abs(config.MigrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for migrations: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", absPath)
	}

	// Add query parameters to connection string for timeouts if they don't exist
	dbURL := config.DatabaseURL
	if config.LockTimeout > 0 {
		if !containsQueryParam(dbURL, "lock_timeout") {
			dbURL = appendQueryParam(dbURL, fmt.Sprintf("lock_timeout=%d", config.LockTimeout))
		}
	}

	if config.StatementTimeout > 0 {
		if !containsQueryParam(dbURL, "statement_timeout") {
			dbURL = appendQueryParam(dbURL, fmt.Sprintf("statement_timeout=%d", config.StatementTimeout))
		}
	}

	// Note: The default_query_exec_mode parameter is specific to pgx driver and not supported by lib/pq
	// We're using lib/pq with golang-migrate so we should not add this parameter

	// Create a new migrate instance
	m, err := migrate.New(
		fmt.Sprintf("file://%s", absPath),
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Set logger
	m.Log = &MigrateLogger{}

	// Run migrations
	startTime := time.Now()
	log.Println("Starting database migrations...")

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Printf("Migrations completed successfully in %v\n", time.Since(startTime))
	return nil
}

// DumpSchema dumps the current database schema to a file
func DumpSchema(dbURL, outputFile string) error {
	if outputFile == "" {
		outputFile = "./schema.sql"
	}

	// Run pg_dump to extract the schema
	cmd := fmt.Sprintf("pg_dump --schema-only --no-owner --no-acl -d %s -f %s", dbURL, outputFile)

	// For security, don't log the actual URL with credentials
	log.Printf("Dumping schema to %s...", outputFile)

	err := execCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to dump schema: %w", err)
	}

	log.Printf("Schema successfully dumped to %s", outputFile)
	return nil
}

// MigrateLogger implements migrate.Logger interface
type MigrateLogger struct{}

// Printf prints migration logs
func (l *MigrateLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Verbose returns whether verbose logging is enabled
func (l *MigrateLogger) Verbose() bool {
	return false
}

// Helper function to execute a shell command
func execCommand(cmdStr string) error {
	parts := strings.Fields(cmdStr)
	cmd := exec.Command(parts[0], parts[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %w - output: %s", err, string(output))
	}

	return nil
}

// Helper functions for URL manipulation
func containsQueryParam(url, param string) bool {
	parts := strings.Split(url, "?")
	if len(parts) < 2 {
		return false
	}

	queryPart := parts[1]
	queryParams := strings.Split(queryPart, "&")

	for _, qp := range queryParams {
		if strings.HasPrefix(qp, param+"=") {
			return true
		}
	}

	return false
}

func appendQueryParam(url, param string) string {
	if url == "" {
		return url
	}

	if strings.HasSuffix(url, "?") {
		return url + param
	}

	if strings.Contains(url, "?") {
		return url + "&" + param
	}

	return url + "?" + param
}
