package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// App
	AppPort string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret     string
	JWTExpireHours string

	// Google OAuth
    GoogleClientID     string
    GoogleClientSecret string
    GoogleRedirectURL  string

	// Gemini
    GeminiAPIKey string

    // Google Cloud
    GoogleSTTCredentials string
}

var App *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	App = &Config{
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "empathify"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpireHours: getEnv("JWT_EXPIRE_HOURS", "24"),
		
        GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
        GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
        GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		
        GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),

        GoogleSTTCredentials: getEnv("GOOGLE_STT_CREDENTIALS", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}