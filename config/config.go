package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	//App
	AppEnv  string
	AppName string

	//User Service
	UserHTTPPort string
	UserGRPCPort string

	//Author Service
	AuthorHTTPPort string
	AuthorGRPCPort string

	//Blog Service
	BlogHTTPPort string
	BlogGRPCPort string

	//Notification Service
	NotificationHTTPPort string
	NotificationGRPCPort string

	// Postgres
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSLMode  string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string

	//JWT
	JWTSecret              string
	JWTAccessTokenExpMin   int
	JWTRefreshTokenExpDays int

	//RabbitMQ
	RabbitHost     string
	RabbitPort     string
	RabbitUser     string
	RabbitPassword string

	RabbitExchange     string
	RabbitExchangeType string
	RabbitBlogRouting  string
	RabbitNotifyQueue  string

	//gRPC
	GRPCTimeoutSec int
	GRPCRetryCount int
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(" .env not found, using system env")
	}

	return &Config{
		//App
		AppEnv:  getEnv("APP_ENV"),
		AppName: getEnv("APP_NAME"),

		//User
		UserHTTPPort: getEnv("USER_SERVICE_HTTP_PORT"),
		UserGRPCPort: getEnv("USER_SERVICE_GRPC_PORT"),

		//Author
		AuthorHTTPPort: getEnv("AUTHOR_SERVICE_HTTP_PORT"),
		AuthorGRPCPort: getEnv("AUTHOR_SERVICE_GRPC_PORT"),

		//Blog
		BlogHTTPPort: getEnv("BLOG_SERVICE_HTTP_PORT"),
		BlogGRPCPort: getEnv("BLOG_SERVICE_GRPC_PORT"),

		//Notification
		NotificationHTTPPort: getEnv("NOTIFICATION_SERVICE_HTTP_PORT"),
		NotificationGRPCPort: getEnv("NOTIFICATION_SERVICE_GRPC_PORT"),

		//Postgres
		PostgresHost:     getEnv("POSTGRES_HOST"),
		PostgresPort:     getEnv("POSTGRES_PORT"),
		PostgresUser:     getEnv("POSTGRES_USER"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD"),
		PostgresDB:       getEnv("POSTGRES_DB"),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE"),

		//Redis
		RedisHost:     getEnv("REDIS_HOST"),
		RedisPort:     getEnv("REDIS_PORT"),
		RedisPassword: getEnv("REDIS_PASSWORD"),
		RedisDB:       getEnv("REDIS_DB"),

		//JWT
		JWTSecret:              getEnv("RABBITMQ_HOST"),
		JWTAccessTokenExpMin:   getEnvAsInt("JWT_ACCESS_TOKEN_EXP_MIN"),
		JWTRefreshTokenExpDays: getEnvAsInt("JWT_REFRESH_TOKEN_EXP_DAY"),

		//RabbitMQ
		RabbitHost:     getEnv("RABBITMQ_HOST"),
		RabbitPort:     getEnv("RABBITMQ_PORT"),
		RabbitUser:     getEnv("RABBITMQ_USER"),
		RabbitPassword: getEnv("RABBITMQ_PASSWORD"),

		RabbitExchange:     getEnv("RABBITMQ_EXCHANGE"),
		RabbitExchangeType: getEnv("RABBITMQ_EXCHANGE_TYPE"),
		RabbitBlogRouting:  getEnv("RABBITMQ_BLOG_CREATED_ROUTING_KEY"),
		RabbitNotifyQueue:  getEnv("RABBITMQ_NOTIFICATION_QUEUE"),

		//gRPC
		GRPCTimeoutSec: getEnvAsInt("GRPC_TIMEOUT_SEC"),
		GRPCRetryCount: getEnvAsInt("GRPC_RETRY_COUNT"),
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing env: %s", key)
	}
	return val
}

func getEnvAsInt(key string) int {
	valStr := getEnv(key)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("Invalid int env: %s", key)
	}
	return val
}
