package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	APIPort        string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiry      time.Duration
	RefreshExpiry  time.Duration
	StoragePath    string
	MigrationsPath string
	AllowedOrigins string
}

func Load() *Config {
	return &Config{
		APIPort:        getEnv("API_PORT", "2020"),
		DatabaseURL:    mustEnv("DATABASE_URL"),
		JWTSecret:      mustEnv("JWT_SECRET"),
		JWTExpiry:      time.Duration(getEnvInt("JWT_EXPIRY_HOURS", 24)) * time.Hour,
		RefreshExpiry:  time.Duration(getEnvInt("REFRESH_TOKEN_EXPIRY_DAYS", 30)) * 24 * time.Hour,
		StoragePath:    getEnv("STORAGE_PATH", "/var/library"),
		MigrationsPath: getEnv("MIGRATIONS_PATH", "./migrations"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:2026"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		slog.Error("required env not set", "key", key)
		os.Exit(1)
	}
	return v
}
