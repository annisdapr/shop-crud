package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	AppPort	  string
	JWTSecret string
}

var (
	config *Config
	once   sync.Once
)

// LoadConfig loads environment variables from .env file and stores them in the Config struct
func LoadConfig() *Config {
	once.Do(func() {

		if os.Getenv("APP_ENV") != "prod" {
			err := godotenv.Load()
			if err != nil {
				log.Fatalf("Error loading .env file: %v", err)
			}
		}

		// Initialize config from environment variables
		config = &Config{
			DBUrl:     getEnv("DB_URL"),
			AppPort:   getEnv("APP_PORT"),
			JWTSecret: getEnv("JWT_SECRET"),
		}
	})
	return config
}

// GetConfig provides a thread-safe way to access the configuration
func GetConfig() *Config {
	return LoadConfig()
}

// Helper function to retrieve environment variables or default values
func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
