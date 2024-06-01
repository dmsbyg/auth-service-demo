package database

import (
	"database/sql"
	"os"
)

var dbInstance *sql.DB

func NewSQLConnection() (db *sql.DB, err error) {
	if dbInstance != nil {
		return db, nil
	}

	db, err = sql.Open("sqlite3", os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}

	dbInstance = db

	return db, nil
}
