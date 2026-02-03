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

The migrations package provides database migration management with versioning and rollback prevention.

### Migration Files

Create SQL files in your migrations directory with the naming convention: `{version}_{description}.sql`

Examples:
```
001_initial_schema.sql
002_add_users_table.sql
003_create_indexes.sql
```

### Running Migrations

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
- **Ordered Execution**: Runs migrations in alphabetical order
- **Error Handling**: Stops on first failure with detailed error messages
- **Logging**: Logs migration progress and status
- **Idempotent**: Safe to run multiple times

### Example Migration File

`001_initial_schema.sql`:
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