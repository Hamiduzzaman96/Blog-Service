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

// Load reads .env and returns Config struct
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	cfg := &Config{
		AppName: os.Getenv("APP_NAME"),
		AppEnv:  os.Getenv("APP_ENV"),
		UserService: ServiceConfig{
			HTTPPort: os.Getenv("USER_SERVICE_HTTP_PORT"),
			GRPCPort: os.Getenv("USER_SERVICE_GRPC_PORT"),
		},
		AuthorService: ServiceConfig{
			HTTPPort: os.Getenv("AUTHOR_SERVICE_HTTP_PORT"),
			GRPCPort: os.Getenv("AUTHOR_SERVICE_GRPC_PORT"),
		},
		BlogService: ServiceConfig{
			HTTPPort: os.Getenv("BLOG_SERVICE_HTTP_PORT"),
			GRPCPort: os.Getenv("BLOG_SERVICE_GRPC_PORT"),
		},
		NotificationService: ServiceConfig{
			HTTPPort: os.Getenv("NOTIFICATION_SERVICE_HTTP_PORT"),
			GRPCPort: os.Getenv("NOTIFICATION_SERVICE_GRPC_PORT"),
		},
		Postgres: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DB:       os.Getenv("POSTGRES_DB"),
			SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret:          os.Getenv("JWT_SECRET"),
			AccessTokenExp:  mustAtoi(os.Getenv("JWT_ACCESS_TOKEN_EXP_MIN")),
			RefreshTokenExp: mustAtoi(os.Getenv("JWT_REFRESH_TOKEN_EXP_DAY")),
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       mustAtoi(os.Getenv("REDIS_DB")),
		},
		RabbitMQ: RabbitMQConfig{
			Host:                  os.Getenv("RABBITMQ_HOST"),
			Port:                  os.Getenv("RABBITMQ_PORT"),
			User:                  os.Getenv("RABBITMQ_USER"),
			Password:              os.Getenv("RABBITMQ_PASSWORD"),
			Exchange:              os.Getenv("RABBITMQ_EXCHANGE"),
			ExchangeType:          os.Getenv("RABBITMQ_EXCHANGE_TYPE"),
			BlogCreatedRoutingKey: os.Getenv("RABBITMQ_BLOG_CREATED_ROUTING_KEY"),
			NotificationQueue:     os.Getenv("RABBITMQ_NOTIFICATION_QUEUE"),
		},
		GRPCTimeoutSec: mustAtoi(os.Getenv("GRPC_TIMEOUT_SEC")),
		GRPCRetryCount: mustAtoi(os.Getenv("GRPC_RETRY_COUNT")),
		LogLevel:       os.Getenv("LOG_LEVEL"),
	}

	return cfg
}

func mustAtoi(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Invalid integer value: %s", val)
	}
	return i
}
