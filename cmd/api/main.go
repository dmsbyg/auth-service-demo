package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dmsbyg/auth-service-demo/config"
	"github.com/dmsbyg/auth-service-demo/database"
	"github.com/dmsbyg/auth-service-demo/internal/auth"
	"github.com/dmsbyg/auth-service-demo/internal/auth/token"
	"github.com/dmsbyg/auth-service-demo/pkg/logger"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config := config.New()
	db, err := database.NewSQLConnection(config)
	if err != nil {
		log.Panic(err)
	}

	l, cleanup, err := logger.NewLogger(config.LoggerConfig)
	if err != nil {
		log.Panicf("cannot start logger: %s", err)
	}
	defer cleanup() //nolint

	jwtMaker, err := token.NewJWTMaker(config.JWTSecret, config.JwtTokenDuration)
	if err != nil {
		log.Panic(err)
	}
	repo := auth.NewRepository(db, &l)
	service := auth.NewService(repo, jwtMaker, &l)
	httpHandler := auth.NewHttpHandler(service)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      httpHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Panicf("cannot start server: %s", err)
	}
}
