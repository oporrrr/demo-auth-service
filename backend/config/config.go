package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort           string
	AuthCenterBaseURL string
	AuthClientID      string
	AuthClientSecret  string
	DatabaseURL       string
	// Role Service
	RoleServiceURL    string
	RoleServiceAPIKey string
	SystemCode        string
}

var Cfg *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	Cfg = &Config{
		AppPort:           getEnv("APP_PORT", "3000"),
		AuthCenterBaseURL: getEnv("AUTH_CENTER_BASE_URL", "https://auth-service-dev.allkons.com"),
		AuthClientID:      getEnv("AUTH_CLIENT_ID", ""),
		AuthClientSecret:  getEnv("AUTH_CLIENT_SECRET", ""),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		RoleServiceURL:    getEnv("ROLE_SERVICE_URL", "http://localhost:3010"),
		RoleServiceAPIKey: getEnv("ROLE_SERVICE_API_KEY", ""),
		SystemCode:        getEnv("SYSTEM_CODE", ""),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
