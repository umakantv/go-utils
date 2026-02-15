package migrations

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/umakantv/go-utils/logger"
)

const migrationFilePattern = `^\d{14}_[a-zA-Z0-9_]+\.sql$`

// Migrate runs database migrations from the specified directory
func Migrate(db *sqlx.DB, migrationsDir string) error {
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

func createMigrationsTable(db *sqlx.DB) error {
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

	re := regexp.MustCompile(migrationFilePattern)

	var migrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			if !re.MatchString(file.Name()) {
				return nil, fmt.Errorf("invalid migration file %s: must match format <UTC timestamp>_<name>.sql where timestamp is 14 digits and name uses alphanum+underscore", file.Name())
			}
			migrations = append(migrations, filepath.Join(dir, file.Name()))
		}
	}

	// Sort migrations by filename
	sort.Strings(migrations)
	return migrations, nil
}

func runMigration(db *sqlx.DB, filePath string) error {
	version := strings.TrimSuffix(filepath.Base(filePath), ".sql")

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = ?)", version).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		logger.Info(fmt.Sprintf("Migration %s already applied, skipping", version))
		return nil
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Running migration: %s", version))
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(string(content))
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Migration %s applied successfully", version))
	return nil
}
