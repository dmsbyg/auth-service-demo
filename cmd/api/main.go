package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dmsbyg/auth-service-demo/database"
	"github.com/dmsbyg/auth-service-demo/internal/auth"
	"github.com/dmsbyg/auth-service-demo/internal/auth/token"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := database.NewSQLConnection()
	if err != nil {
		log.Panic(err)
	}

	tokenDuration, err := time.ParseDuration(os.Getenv("JWT_TOKEN_DURATION"))
	if err != nil {
		log.Panic(err)
	}

	jwtMaker, err := token.NewJWTMaker(os.Getenv("JWT_SECRET"), tokenDuration)
	if err != nil {
		log.Panic(err)
	}
	repo := auth.NewRepository(db)
	service := auth.NewService(repo, jwtMaker)
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
