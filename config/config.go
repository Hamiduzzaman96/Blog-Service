package config

import (
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
	AccessTokenExp  int // minutes
	RefreshTokenExp int // days
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
	_ = godotenv.Load() // ignore error, env variables may be used

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
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DB:       getEnv("POSTGRES_DB", "blog-service"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "ramim12345ramim12345"),
			AccessTokenExp:  getEnvAsInt("JWT_ACCESS_TOKEN_EXP_MIN", 15),
			RefreshTokenExp: getEnvAsInt("JWT_REFRESH_TOKEN_EXP_DAY", 7),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		RabbitMQ: RabbitMQConfig{
			Host:                  getEnv("RABBITMQ_HOST", "localhost"),
			Port:                  getEnv("RABBITMQ_PORT", "5672"),
			User:                  getEnv("RABBITMQ_USER", "guest"),
			Password:              getEnv("RABBITMQ_PASSWORD", "guest"),
			Exchange:              getEnv("RABBITMQ_EXCHANGE", "blog.events"),
			ExchangeType:          getEnv("RABBITMQ_EXCHANGE_TYPE", "topic"),
			BlogCreatedRoutingKey: getEnv("RABBITMQ_BLOG_CREATED_ROUTING_KEY", "blog.created"),
			NotificationQueue:     getEnv("RABBITMQ_NOTIFICATION_QUEUE", "notification.queue"),
		},
		GRPCTimeoutSec: getEnvAsInt("GRPC_TIMEOUT_SEC", 3),
		GRPCRetryCount: getEnvAsInt("GRPC_RETRY_COUNT", 3),
		LogLevel:       getEnv("LOG_LEVEL", "debug"),
	}

	return cfg
}

// helpers

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
