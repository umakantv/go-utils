# Database Module

The database module provides utilities for connecting to databases and managing migrations using [sqlx](https://github.com/jmoiron/sqlx).

## Supported Databases

- **MySQL**: Requires `github.com/go-sql-driver/mysql` driver
- **PostgreSQL**: Requires `github.com/lib/pq` driver
- **SQLite**: Requires `github.com/mattn/go-sqlite3` driver

Make sure to import the appropriate driver in your main application:

```go
import (
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
    _ "github.com/mattn/go-sqlite3"
)
```

## Configuration

Define your database configuration using the `DatabaseConfig` struct:

```go
type DatabaseConfig struct {
    DRIVER   string
    HOST     string
    PORT     string
    USER     string
    PASSWORD string
    DB       string
}
```

## Connecting to Database

Use `GetDBConnection` to establish a database connection. The function supports MySQL, PostgreSQL, and SQLite.

### MySQL
```go
config := DatabaseConfig{
    DRIVER:   "mysql",
    HOST:     "localhost",
    PORT:     "3306",
    USER:     "user",
    PASSWORD: "password",
    DB:       "mydb",
}

db := GetDBConnection(config)
defer db.Close()
```

### PostgreSQL
```go
config := DatabaseConfig{
    DRIVER:   "postgres",
    HOST:     "localhost",
    PORT:     "5432",
    USER:     "user",
    PASSWORD: "password",
    DB:       "mydb",
}

db := GetDBConnection(config)
defer db.Close()
```

### SQLite (File-based)
```go
config := DatabaseConfig{
    DRIVER: "sqlite3",
    DB:     "./mydb.sqlite", // File path
}

db := GetDBConnection(config)
defer db.Close()
```

### SQLite (In-Memory)
```go
config := DatabaseConfig{
    DRIVER: "sqlite3",
    DB:     ":memory:", // Or leave empty for in-memory
}

db := GetDBConnection(config)
defer db.Close()
```

## Connection Pool Settings

All connections are configured with:
- Max lifetime: 3 minutes
- Max open connections: 10
- Max idle connections: 10

These can be adjusted by modifying the `GetDBConnection` function for your specific needs.

## Migrations

The migrations package provides database migration management with versioning, rollback prevention, and enforced filename patterns for ordering.

### Migration Files

All migration files must follow the exact format: `<UTC timestamp>_<migration_name_with_alphanum_and_underscore>.sql`
where timestamp is 14 digits (e.g. YYYYMMDDHHMMSS from UTC time) to ensure sortable order.
Invalid files cause immediate error before any migrations run.

Examples:
```
20230101120000_initial_schema.sql
20230101120001_add_users_table.sql
20230101120002_create_indexes.sql
```

### Creating Migrations

Run from your project root (creates .sql file under your ./migrations dir; uses UTC timestamp prefix and validates format/name):
```
go run github.com/umakantv/go-utils/db/migrations/create-migration.go --name create_user_table --dir ./migrations
```
(If developing go-utils locally, use `go run db/migrations/create-migration.go ...` instead.)
The script errors immediately for invalid names (e.g. spaces, special chars) or format violations; generated files go under the specified dir (with .sql extension).

### CLI for Running Migrations

For command-line migrations (reuses library logic to skip already-applied files based on schema_migrations table, validates filenames, etc.):
Place a `.env` file in the current dir with DB config (keys match former flags):
```
DRIVER=sqlite3
DB=./mydb.sqlite
# Optional: HOST=localhost
# PORT=3306
# USER=root
# PASSWORD=pass
```

Then:
```
# Example (run from project root)
go run github.com/umakantv/go-utils/db/migrations/migrate.go --dir ./migrations
```
(For local go-utils dev: replace with `go run db/migrations/migrate.go ...`.) The --dir flag specifies where to find .sql files.

### Running Migrations (in Go code)

```go
import "github.com/umakantv/go-utils/db/migrations"

// Run all pending migrations
err := migrations.Migrate(db, "./migrations")
if err != nil {
    logger.Error("Migration failed", logger.Error(err))
}
```

### Migration Tracking

The system creates a `schema_migrations` table to track applied migrations:

```sql
CREATE TABLE schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Features

- **Version Control**: Tracks applied migrations to prevent re-execution
- **Ordered Execution**: Runs migrations in alphabetical order by UTC timestamp prefix
- **Filename Validation**: Enforces `<timestamp>_<name>.sql` pattern; errors immediately on mismatch
- **Transactions**: Each migration is wrapped in a DB transaction (rollback on failure)
- **Error Handling**: Stops on first failure with detailed error messages
- **Logging**: Logs migration progress and status
- **Idempotent**: Safe to run multiple times

### Example Migration File

`20230101120000_initial_schema.sql`:
```sql
-- Create users table
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create posts table
CREATE TABLE posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### Best Practices

1. Use descriptive migration names
2. Test migrations on a copy of production data first
3. Keep migrations small and focused
4. Never modify existing migration files
5. Use transactions for complex migrations (if supported by your DB)
6. Version control your migration files
7. Run migrations during application startup