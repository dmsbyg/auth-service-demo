package database

import (
	"os"

	"github.com/jmoiron/sqlx"
)

var dbInstance *sqlx.DB

func NewSQLConnection() (db *sqlx.DB, err error) {
	if dbInstance != nil {
		return db, nil
	}

	db, err = sqlx.Open("sqlite3", os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}

	dbInstance = db

	return db, nil
}
