package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

// private helper function
func getEnvWithDefault(key, defaultValue string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Printf("getting key: %s", key)
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func GetDBConfig() DBConfig {
	const (
		defaultUser     = "postgres"
		defaultPassword = "postgres"
		defaultDatabase = "database"
		defaultHost     = "localhost"
		defaultPort     = "5432"
	)

	return DBConfig{
		User:     getEnvWithDefault("DB_USER", defaultUser),
		Password: getEnvWithDefault("DB_PASSWORD", defaultPassword),
		Database: getEnvWithDefault("DB_DATABASE", defaultDatabase),
		Host:     getEnvWithDefault("DB_HOST", defaultHost),
		Port:     getEnvWithDefault("DB_PORT", defaultPort),
	}
}
