package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	BACKEND_PORT                string
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
	KestraURL           string
	KestraUsername      string
	KestraPassword      string
	BackendURL          string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found")
	}

	return &Config{
		BACKEND_PORT:                getEnv("PORT", "8081"),
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
		KestraURL:           getEnv("KESTRA_URL", "http://localhost:8080"),
		KestraUsername:      getEnv("KESTRA_USERNAME", ""),
		KestraPassword:      getEnv("KESTRA_PASSWORD", ""),
		BackendURL:          getEnv("BACKEND_URL", "http://host.docker.internal:8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
