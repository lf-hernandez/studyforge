package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	ServerPort       string
	ServerHost       string
	DatabasePath     string
	UploadDir        string
	MaxFileSize      int64
	SessionTimeout   int
	HuggingFaceKey   string
	HuggingFaceURL   string
	LogLevel         string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		ServerHost:       getEnv("SERVER_HOST", "localhost"),
		DatabasePath:     getEnv("DATABASE_PATH", "./data/studyforge.db"),
		UploadDir:        getEnv("UPLOAD_DIR", "./uploads"),
		MaxFileSize:      getEnvInt64("MAX_FILE_SIZE", 52428800), // 50MB
		SessionTimeout:   getEnvInt("SESSION_TIMEOUT", 86400),     // 24 hours
		HuggingFaceKey:   getEnv("HUGGINGFACE_API_KEY", ""),
		HuggingFaceURL:   getEnv("HUGGINGFACE_API_URL", "https://api-inference.huggingface.co/models"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv retrieves environment variable or returns default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt retrieves integer environment variable or returns default
func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid integer for %s, using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}

// getEnvInt64 retrieves int64 environment variable or returns default
func getEnvInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		log.Printf("Invalid int64 for %s, using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}
