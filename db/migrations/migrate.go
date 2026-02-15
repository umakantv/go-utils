package migrations

import (
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/umakantv/go-utils/db"
	"github.com/umakantv/go-utils/logger"
)

func RunMigrations() {
	dirFlag := flag.String("dir", "", "Directory containing migration .sql files")
	flag.Parse()

	if *dirFlag == "" {
		fmt.Println("Usage: go run migrate.go --dir <migrations-dir>")
		os.Exit(1)
	}

	config := db.DatabaseConfig{}
	f, _ := os.ReadFile(".env")
	for _, line := range strings.Split(string(f), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			switch key {
			case "DRIVER":
				config.DRIVER = val
			case "HOST":
				config.HOST = val
			case "PORT":
				config.PORT = val
			case "USER":
				config.USER = val
			case "PASSWORD":
				config.PASSWORD = val
			case "DB":
				config.DB = val
			}
		}
	}
	if config.DRIVER == "" || config.DB == "" {
		fmt.Println("Error: .env missing required keys (DRIVER, DB) or file not found")
		os.Exit(1)
	}

	logger.Init(logger.LoggerConfig{
		CallerKey:  "file",
		TimeKey:    "timestamp",
		CallerSkip: 1,
	})

	sqlxDB := db.GetDBConnection(config)
	defer sqlxDB.Close()

	if err := Migrate(sqlxDB, *dirFlag); err != nil {
		fmt.Printf("Migration failed: %v\n", err)
		os.Exit(1)
	}
}
