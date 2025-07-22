package config

import (
	"os"
)

type Config struct {
	JWTSecret   string
	DBUser      string
	DBPassword  string
	DBHost      string
	DBPort      string
	DBName      string
	DBURL       string
	LlamaAPIKey string
	LLMURL      string
	REDIS_URL   string
}

var AppConfig Config

func init() {
	// Load environment variables from .env for local development
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("Error loading .env file, relying on system environment variables")
	// }

	// Assign environment variables
	AppConfig = Config{
		JWTSecret:   getEnv("JWT_SECRET", "default_jwt_secret"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBName:      getEnv("DB_NAME", "readeaze"),
		DBURL:       getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/readeaze"),
		LlamaAPIKey: getEnv("LLAMA_API_KEY", "default_llama_key"),
		LLMURL:      getEnv("LLM_URL", "http://localhost:8000"),
		REDIS_URL:   getEnv("REDIS_URL", "redis://localhost:6379/0"),
	}
}

// Helper to fetch environment variables
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
