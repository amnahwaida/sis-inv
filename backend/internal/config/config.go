package config

import (
	"os"
	"time"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	JWTSecret      string
	JWTAccessExp   time.Duration
	JWTRefreshExp  time.Duration
	ServerPort     string
	GinMode        string
}

func Load() *Config {
	accessExp, _ := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "24h"))
	refreshExp, _ := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))

	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "sekolah_admin"),
		DBPassword:    getEnv("DB_PASSWORD", "sekolah_secret_2024"),
		DBName:        getEnv("DB_NAME", "inventaris_db"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		JWTSecret:     getEnv("JWT_SECRET", "sis-inv-dev-jwt-secret-key-2024"),
		JWTAccessExp:  accessExp,
		JWTRefreshExp: refreshExp,
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		GinMode:       getEnv("GIN_MODE", "debug"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
