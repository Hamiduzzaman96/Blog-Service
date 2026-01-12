package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ServiceConfig struct {
	HTTPPort string
	GRPCPort string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret          string
	AccessTokenExp  int
	RefreshTokenExp int
}

type RabbitMQConfig struct {
	Host                  string
	Port                  string
	User                  string
	Password              string
	Exchange              string
	ExchangeType          string
	BlogCreatedRoutingKey string
	NotificationQueue     string
}

type Config struct {
	AppName             string
	AppEnv              string
	UserService         ServiceConfig
	AuthorService       ServiceConfig
	BlogService         ServiceConfig
	NotificationService ServiceConfig
	Postgres            PostgresConfig
	Redis               RedisConfig
	JWT                 JWTConfig
	RabbitMQ            RabbitMQConfig
	GRPCTimeoutSec      int
	GRPCRetryCount      int
	LogLevel            string
}

// Load reads .env and environment variables
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using environment variables")
	}

	cfg := &Config{
		AppName: os.Getenv("APP_NAME"),
		AppEnv:  getEnv("APP_ENV", "development"),

		UserService: ServiceConfig{
			HTTPPort: getEnv("USER_SERVICE_HTTP_PORT", ":8001"),
			GRPCPort: getEnv("USER_SERVICE_GRPC_PORT", ":50051"),
		},
		AuthorService: ServiceConfig{
			HTTPPort: getEnv("AUTHOR_SERVICE_HTTP_PORT", ":8002"),
			GRPCPort: getEnv("AUTHOR_SERVICE_GRPC_PORT", ":50052"),
		},
		BlogService: ServiceConfig{
			HTTPPort: getEnv("BLOG_SERVICE_HTTP_PORT", ":8003"),
			GRPCPort: getEnv("BLOG_SERVICE_GRPC_PORT", ":50053"),
		},
		NotificationService: ServiceConfig{
			HTTPPort: getEnv("NOTIFICATION_SERVICE_HTTP_PORT", ":8004"),
			GRPCPort: getEnv("NOTIFICATION_SERVICE_GRPC_PORT", ":50054"),
		},

		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", ""),
			Port:     getEnv("POSTGRES_PORT", ""),
			User:     getEnv("POSTGRES_USER", ""),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			DB:       getEnv("POSTGRES_DB", ""),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", ""),
			Port:     getEnv("REDIS_PORT", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},

		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", ""),
			AccessTokenExp:  getEnvAsInt("JWT_ACCESS_TOKEN_EXP_MIN", 15),
			RefreshTokenExp: getEnvAsInt("JWT_REFRESH_TOKEN_EXP_DAY", 7),
		},

		RabbitMQ: RabbitMQConfig{
			Host:                  getEnv("RABBITMQ_HOST", ""),
			Port:                  getEnv("RABBITMQ_PORT", ""),
			User:                  getEnv("RABBITMQ_USER", ""),
			Password:              getEnv("RABBITMQ_PASSWORD", ""),
			Exchange:              getEnv("RABBITMQ_EXCHANGE", ""),
			ExchangeType:          getEnv("RABBITMQ_EXCHANGE_TYPE", ""),
			BlogCreatedRoutingKey: getEnv("RABBITMQ_BLOG_CREATED_ROUTING_KEY", ""),
			NotificationQueue:     getEnv("RABBITMQ_NOTIFICATION_QUEUE", ""),
		},

		GRPCTimeoutSec: getEnvAsInt("GRPC_TIMEOUT_SEC", 3),
		GRPCRetryCount: getEnvAsInt("GRPC_RETRY_COUNT", 3),
		LogLevel:       getEnv("LOG_LEVEL", "debug"),
	}

	return cfg
}

// helpers
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
