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
		DBUser:      getEnv("DB_USER", "default_user"),
		DBPassword:  getEnv("DB_PASSWORD", "default_password"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBName:      getEnv("DB_NAME", "default_db"),
		DBURL:       getEnv("DB_URL", "postgres://default_user:default_password@localhost:5432/default_db"),
		LlamaAPIKey: getEnv("LLAMA_API_KEY", "default_llama_key"),
		LLMURL:      getEnv("LLM_URL", "http://localhost:8000"),
	}
}

// Helper to fetch environment variables
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
