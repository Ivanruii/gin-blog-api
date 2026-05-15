package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBPath     string
	GinMode    string
	AppVersion string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBPath:     getEnv("DB_PATH", "blog.db"),
		GinMode:    getEnv("GIN_MODE", "debug"),
		AppVersion: getEnv("APP_VERSION", "0.1.0"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
