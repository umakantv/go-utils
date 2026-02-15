package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func CreateMigration(nameFlag *string, dirFlag *string) {

	if *nameFlag == "" {
		fmt.Println("Usage: go run create-migration.go --name <name> [--dir <migrations-dir>]")
		os.Exit(1)
	}

	name := strings.TrimSpace(*nameFlag)
	if name != *nameFlag {
		fmt.Printf("Error: name contains leading/trailing spaces: %s\n", *nameFlag)
		os.Exit(1)
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(name) {
		fmt.Printf("Error: invalid name '%s' - must contain only alphanumeric chars and underscores (no spaces or special chars)\n", name)
		os.Exit(1)
	}

	timestamp := time.Now().UTC().Format("20060102150405")
	filename := timestamp + "_" + name + ".sql"
	fullPath := filepath.Join(*dirFlag, filename)

	if !regexp.MustCompile(migrationFilePattern).MatchString(filename) {
		fmt.Printf("Error: generated filename '%s' does not meet required format\n", filename)
		os.Exit(1)
	}

	if _, err := os.Stat(fullPath); err == nil {
		fmt.Printf("Error: file already exists: %s\n", fullPath)
		os.Exit(1)
	}

	if err := os.MkdirAll(*dirFlag, 0755); err != nil {
		fmt.Printf("Error creating dir: %v\n", err)
		os.Exit(1)
	}

	content := fmt.Sprintf(`-- Migration: %s
-- Generated: %s UTC

-- Add your SQL migration here
`, name, timestamp)
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created migration file: %s\n", fullPath)
}
