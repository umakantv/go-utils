package db

import (
	"fmt"
	"time"

	"github.com/umakantv/go-utils/logger"

	"github.com/jmoiron/sqlx"
)

func GetDBConnection(dbConfig DatabaseConfig) *sqlx.DB {
	var dsn string

	switch dbConfig.DRIVER {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			dbConfig.USER, dbConfig.PASSWORD, dbConfig.HOST, dbConfig.PORT, dbConfig.DB)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbConfig.HOST, dbConfig.PORT, dbConfig.USER, dbConfig.PASSWORD, dbConfig.DB)
	case "sqlite3":
		// For SQLite, DB field contains the file path or ":memory:"
		dsn = dbConfig.DB
		if dsn == "" {
			dsn = ":memory:"
		}
	default:
		dsn = fmt.Sprintf("%s:%s@/%s", dbConfig.USER, dbConfig.PASSWORD, dbConfig.DB)
	}

	db, err := sqlx.Open(dbConfig.DRIVER, dsn)
	if err != nil {
		logger.Error("Error in opening a DB connection " + err.Error())
		panic(err) // or return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error("Error in ping to DB connection " + err.Error())
		panic(err) // or return nil, err
	}

	// Connection pool settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
