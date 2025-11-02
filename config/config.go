package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL                     string
	Port                            string
	OpenAIAPIKey                    string
	TranslationCacheTTL             time.Duration
	TranslationCacheCleanupInterval time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	config := &Config{
		DatabaseURL:                     getEnv("DB_URL"),
		Port:                            getEnvWithDefault("PORT", "8080"),
		OpenAIAPIKey:                    os.Getenv("OPENAI_API_KEY"), // Optional
		TranslationCacheTTL:             getDurationEnvWithDefault("TRANSLATION_CACHE_TTL", 1*time.Hour),
		TranslationCacheCleanupInterval: getDurationEnvWithDefault("TRANSLATION_CACHE_CLEANUP_INTERVAL", 10*time.Minute),
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

func getDurationEnvWithDefault(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Try to parse as integer seconds
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}

	// Try to parse as duration string (e.g., "1h", "30m")
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}

	log.Printf("Invalid duration for %s: %s, using default: %s", key, value, defaultValue)
	return defaultValue
}
