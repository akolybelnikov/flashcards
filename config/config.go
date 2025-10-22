package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	config := &Config{
		DatabaseURL: getEnv("DB_URL"),
		Port:        getEnvWithDefault("PORT", "8080"),
	}

	return config
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Required environment variable not set: " + key)
	}
	return value
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
