package config

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	DBHost  string
	DBPort  string
	DBUser  string
	DBPass  string
	DBName  string
}

func Load() (*Config, error) {
	// Use Overload so .env values override existing OS env vars
	_ = godotenv.Overload()

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		DBHost:  getEnv("DB_HOST", ""),
		DBPort:  getEnv("DB_PORT", "3306"),
		DBUser:  getEnv("DB_USER", ""),
		DBPass:  getEnv("DB_PASS", ""),
		DBName:  getEnv("DB_NAME", ""),
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" {
		return nil, errors.New("missing required DB env vars: DB_HOST, DB_USER, DB_NAME")
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}
