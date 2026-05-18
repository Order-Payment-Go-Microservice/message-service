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
	NatsURL                 string
	RedisURL                string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:                    getEnv("PORT", "9004"),
		GRPCPort:                getEnv("GRPC_PORT", "50051"),
		DBHost:                  getEnv("DB_HOST", "localhost"),
		DBPort:                  getEnv("DB_PORT", "5432"),
		DBUser:                  getEnv("DB_USER", "postgres"),
		DBPassword:              getEnv("DB_PASSWORD", "postgres"),
		DBName:                  getEnv("DB_NAME", "messages_db"),
		NotificationServiceAddr: getEnv("NOTIFICATION_SERVICE_ADDR", "localhost:50052"),
		NatsURL:                 getEnv("NATS_URL", "nats://localhost:4222"),
		RedisURL:                getEnv("REDIS_URL", "localhost:6379"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
