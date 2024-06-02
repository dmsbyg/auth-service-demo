package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Panic("DB_URL is not set, plese check your environment variables")
	}
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close() //nolint
	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Panic(err)
	}

	fSrc, err := (&file.File{}).Open("./database/migrations")
	if err != nil {
		log.Panic(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
	if err != nil {
		log.Panic(err)
	}
	log.Println("migrations started")

	if err := m.Up(); err != nil {
		log.Println("panic here")
		log.Panic(err)
	}

	log.Println("migrations success")
}
