package database

import (
	"github.com/dmsbyg/auth-service-demo/config"
	"github.com/jmoiron/sqlx"
)

var dbInstance *sqlx.DB

func NewSQLConnection(config *config.Config) (db *sqlx.DB, err error) {
	if dbInstance != nil {
		return db, nil
	}

	db, err = sqlx.Open("sqlite3", config.DBUrl)
	if err != nil {
		return nil, err
	}

	dbInstance = db

	return db, nil
}
