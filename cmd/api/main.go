package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	log.Println("Initializing... load configuration...")
	config := config.New()
	db, err := database.NewSQLConnection(config)
	if err != nil {
		log.Panic(err)
	}

	logger, cleanup, err := logger.NewLogger(config.LoggerConfig)
	if err != nil {
		log.Panicf("cannot start logger: %s", err)
	}
	defer cleanup() //nolint

	jwtMaker, err := token.NewJWTMaker(config.JWTSecret, config.JwtTokenDuration)
	if err != nil {
		log.Panic(err)
	}
	repo := auth.NewRepository(db, &logger)
	service := auth.NewService(repo, jwtMaker, &logger)
	httpHandler := auth.NewHTTPHandler(service)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      httpHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		logger.Infof("Starting server on port :%d", config.Port)
		log.Printf("Starting server on port :%d \n", config.Port)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicf("cannot start server: %s", err)
		}
	}()

	<-ctx.Done()
	log.Println("gracefully shut down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Panicf("cannot shutdown server: %s", err)
	}
}
