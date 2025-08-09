package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath               string
	JWTSecret           string
	JWTExpiry           time.Duration
	RefreshTokenExpiry  time.Duration
	Port                string
	CORSOrigins         string
	LogLevel            string
	RateLimit           int
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT", "100"))
	
	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "15m"))
	refreshExpiry, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h"))

	return &Config{
		DBPath:             getEnv("DB_PATH", "./analytics.db"),
		JWTSecret:          getEnv("JWT_SECRET", "your-jwt-secret-key"),
		JWTExpiry:          jwtExpiry,
		RefreshTokenExpiry: refreshExpiry,
		Port:               getEnv("PORT", "8080"),
		CORSOrigins:        getEnv("CORS_ORIGINS", "http://localhost:3000"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		RateLimit:          rateLimit,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}