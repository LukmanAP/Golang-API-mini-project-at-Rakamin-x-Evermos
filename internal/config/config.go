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
	// Address/EMSIFA configuration
	EMSIFABase       string
	HTTPTimeoutMS    int
	HTTPRetry        int
	CacheTTLSeconds  int
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
		// Defaults for EMSIFA-based address service
		EMSIFABase:       getEnv("EMSIFA_BASE", "https://www.emsifa.com/api-wilayah-indonesia/api"),
		HTTPTimeoutMS:    getEnvInt("HTTP_TIMEOUT_MS", 5000),
		HTTPRetry:        getEnvInt("HTTP_RETRY", 2),
		CacheTTLSeconds:  getEnvInt("CACHE_TTL_SECONDS", 86400),
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
