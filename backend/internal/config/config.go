package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	DBHost              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBPort              string
	GithubClientID      string
	GithubClientSecret  string
	GithubWebhookSecret string
	GithubRedirectURI   string
	FrontendURL         string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		Port:                getEnv("PORT", "8080"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "postgres"),
		DBName:              getEnv("DB_NAME", "devplus"),
		DBPort:              getEnv("DB_PORT", "5432"),
		GithubClientID:      getEnv("GITHUB_CLIENT_ID", "Ov23li25TeWbIlsfacBJ"),
		GithubClientSecret:  getEnv("GITHUB_CLIENT_SECRET", ""),
		GithubWebhookSecret: getEnv("GITHUB_WEBHOOK_SECRET", ""),
		GithubRedirectURI:   getEnv("GITHUB_REDIRECT_URI", ""),
		FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:3000/dashboard"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
