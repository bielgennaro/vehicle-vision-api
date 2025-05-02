package configs

import (
	"os"
)

type Config struct {
	Port             string
	Environment      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
}

func LoadConfig() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		Environment:      getEnv("ENVIRONMENT", "development"),
		DatabaseHost:     getEnv("DATABASE_URL", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "5432"),
		DatabaseName:     getEnv("DATABASE_NAME", "vehicle_vision"),
		DatabaseUser:     getEnv("DATABASE_USER", "postgres"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", "99831"),
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
