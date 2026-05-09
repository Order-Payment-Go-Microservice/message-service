package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    string
	GRPCPort                string
	DBHost                  string
	DBPort                  string
	DBUser                  string
	DBPassword              string
	DBName                  string
	NotificationServiceAddr string
	WSPort                  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
	}

	return &Config{
		Port:                    getEnv("PORT", "9004"),
		GRPCPort:                getEnv("GRPC_PORT", "50051"),
		DBHost:                  getEnv("DB_HOST", "localhost"),
		DBPort:                  getEnv("DB_PORT", "5432"),
		DBUser:                  getEnv("DB_USER", "postgres"),
		DBPassword:              getEnv("DB_PASSWORD", "postgres"),
		DBName:                  getEnv("DB_NAME", "messages_db"),
		NotificationServiceAddr: getEnv("NOTIFICATION_SERVICE_ADDR", "localhost:50052"),
		WSPort:                  getEnv("WS_PORT", "9005"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
