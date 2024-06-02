package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dmsbyg/auth-service-demo/pkg/logger"
)

type Config struct {
	Port             int           `mapstructure:"PORT"`
	AppEnv           string        `mapstructure:"APP_ENV"`
	DBUrl            string        `mapstructure:"DB_URL"`
	JWTSecret        string        `mapstructure:"JWT_SECRET"`
	JwtTokenDuration time.Duration `mapstructure:"JWT_TOKEN_DURATION"`

	LoggerConfig logger.Config
}

func New() *Config {
	cfg := newDefaultConfig()

	if val := getenv("APP_ENV"); val != "" {
		cfg.AppEnv = val
	}
	if val := getenv("PORT"); val != "" {
		port, err := strconv.Atoi(val)
		if err == nil {
			cfg.Port = port
		}
	}
	if val := getenv("DB_URL"); val != "" {
		cfg.DBUrl = val
	}
	if val := getenv("JWT_SECRET"); val != "" {
		cfg.JWTSecret = val
	}
	if val := getenv("JWT_TOKEN_DURATION"); val != "" {
		duration, err := time.ParseDuration(val)
		if err == nil {
			cfg.JwtTokenDuration = duration
		}
	}
	if val := getenv("LOG_LEVEL"); val != "" {
		cfg.LoggerConfig.LogLevel = val
	}
	if val := getenv("LOG_ENCODING"); val != "" {
		cfg.LoggerConfig.LogEncoding = val
	}
	if val := getenv("LOG_OUTPUT"); val != "" {
		cfg.LoggerConfig.LogOutput = val
	}
	if val := getenv("LOG_ERROR_OUTPUT"); val != "" {
		cfg.LoggerConfig.LogErrorOutput = val
	}

	return cfg
}

func newDefaultConfig() *Config {
	return &Config{
		Port:             8080,
		AppEnv:           "development",
		DBUrl:            "./test.db",
		JWTSecret:        "jwtsecretwhichhasenoughcharacters",
		JwtTokenDuration: 1 * time.Hour,
		LoggerConfig: logger.Config{
			LogLevel:       "error",
			LogEncoding:    "console",
			LogOutput:      "stdout",
			LogErrorOutput: "stderr",
		},
	}
}

func getenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("cannot get variable for: %s, use default value\n", key)
	}
	return val
}
