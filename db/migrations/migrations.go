package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/umakantv/go-utils/logger"
)

// Migrate runs database migrations from the specified directory
func Migrate(db *sql.DB, migrationsDir string) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	files, err := getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Run migrations
	for _, file := range files {
		if err := runMigration(db, file); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file, err)
		}
	}

	logger.Info("All migrations completed successfully")
	return nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func getMigrationFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrations = append(migrations, filepath.Join(dir, file.Name()))
		}
	}

	// Sort migrations by filename
	sort.Strings(migrations)
	return migrations, nil
}

func runMigration(db *sql.DB, filePath string) error {
	// Extract version from filename (e.g., "001_initial.sql" -> "001_initial")
	version := strings.TrimSuffix(filepath.Base(filePath), ".sql")

	// Check if migration already applied
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = ?)", version).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		logger.Info(fmt.Sprintf("Migration %s already applied, skipping", version))
		return nil
	}

	// Read migration file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Execute migration
	logger.Info(fmt.Sprintf("Running migration: %s", version))
	_, err = db.Exec(string(content))
	if err != nil {
		return err
	}

	// Record migration as applied
	_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Migration %s applied successfully", version))
	return nil
}
