package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPass          string
	DBName          string
	JWTSecret       string
	JWTExpiryDays   int
	UploadDirProduct string
	BaseFileURL      string
}

func Load() (*Config, error) {
	// Use Overload so .env values override existing OS env vars
	_ = godotenv.Overload()

	cfg := &Config{
		AppPort:         getEnv("APP_PORT", ""),
		DBHost:          getEnv("DB_HOST", ""),
		DBPort:          getEnv("DB_PORT", ""),
		DBUser:          getEnv("DB_USER", ""),
		DBPass:          getEnv("DB_PASS", ""),
		DBName:          getEnv("DB_NAME", ""),
		JWTSecret:       getEnv("JWT_SECRET", ""),
		JWTExpiryDays:   getEnvInt("JWT_EXP_DAYS", 7),
		UploadDirProduct: getEnv("UPLOAD_DIR_PRODUCT", "uploads/products"),
		BaseFileURL:      getEnv("BASE_FILE_URL", ""),
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" {
		return nil, errors.New("missing required DB env vars: DB_HOST, DB_USER, DB_NAME")
	}

	if strings.TrimSpace(cfg.JWTSecret) == "" {
		return nil, errors.New("missing required JWT env var: JWT_SECRET")
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

func getEnvInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
