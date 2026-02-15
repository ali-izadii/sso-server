package config

import (
	"os"
	"time"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	JWTExpiry          time.Duration
	RefreshTokenExpiry time.Duration
	AuthCodeExpiry     time.Duration
}

func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "sqlite://sso.db"),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpiry:          getDuration("JWT_EXPIRY", time.Hour),
		RefreshTokenExpiry: getDuration("REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
		AuthCodeExpiry:     getDuration("AUTH_CODE_EXPIRY", 10*time.Minute),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
