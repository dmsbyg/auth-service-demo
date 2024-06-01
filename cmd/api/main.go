package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dmsbyg/auth-service-demo/database"
	"github.com/dmsbyg/auth-service-demo/internal/auth"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

func main() { //nolint
	db, err := database.NewSQLConnection()
	if err != nil {
		log.Panic(err)
	}

	repo := auth.NewRepository(db)
	service := auth.NewService(repo)
	httpHandler := auth.NewHttpHandler(service)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Panicf("incorrect port config: %s", err)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: httpHandler,
	}

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Panicf("cannot start server: %s", err)
	}
}
