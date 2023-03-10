package db

import (
	"fmt"
	"time"

	"github.com/umakantv/go-utils/logger"

	"github.com/jmoiron/sqlx"
)

func GetDBConnection(dbConfig DatabaseConfig) *sqlx.DB {

	fmt.Println(dbConfig.DRIVER)
	// Use process env variables here instead for this
	db, err := sqlx.Open(dbConfig.DRIVER, fmt.Sprintf("%v:%v@/%v", dbConfig.USER, dbConfig.PASSWORD, dbConfig.DB))
	if err != nil {
		logger.Error("Error in opening a DB connection " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		logger.Error("Error in ping to DB connection " + err.Error())
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
