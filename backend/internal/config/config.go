package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Load loads environment variables from a .env file
func Load(path string) error {
	// Try to load .env file but don't fail if it doesn't exist

	return godotenv.Load(path)
}

// GetEnv retrieves an environment variable with a default value
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
