package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment   string
	Port          string
	MongoURI      string
	MongoDatabase string
	RelayAPIKey   string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil &&
		!errors.Is(err, os.ErrNotExist) {
		return Config{}, errors.New("failed to load configuration")
	}

	cfg := Config{
		Environment:   getEnv("APP_ENV", "development"),
		Port:          getEnv("PORT", "8080"),
		MongoURI:      strings.TrimSpace(os.Getenv("MONGODB_URI")),
		MongoDatabase: getEnv("MONGODB_DATABASE", "relay"),
		RelayAPIKey:   strings.TrimSpace(os.Getenv("RELAY_API_KEY")),
	}

	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, errors.New(
			"PORT must be between 1 and 65535",
		)
	}

	if cfg.MongoURI == "" {
		return Config{}, errors.New("MONGODB_URI is required")
	}

	if cfg.MongoDatabase == "" {
		return Config{}, errors.New("MONGODB_DATABASE is required")
	}

	if len(cfg.RelayAPIKey) < 32 {
		return Config{}, errors.New(
			"RELAY_API_KEY must contain at least 32 characters",
		)
	}

	if cfg.Environment != "development" &&
		cfg.Environment != "test" &&
		cfg.Environment != "production" {
		return Config{}, errors.New("invalid APP_ENV value")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
