package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hamiduzzaman96/Blog-Service/config"
	grpcHandler "github.com/Hamiduzzaman96/Blog-Service/internal/handler/grpc"
	httpHandler "github.com/Hamiduzzaman96/Blog-Service/internal/handler/http"
	"github.com/Hamiduzzaman96/Blog-Service/internal/repository"
	"github.com/Hamiduzzaman96/Blog-Service/internal/usecase"
	"github.com/Hamiduzzaman96/Blog-Service/pkg/rabbitmq"
	"github.com/Hamiduzzaman96/Blog-Service/proto/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	// Context + Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Postgres connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DB,
		cfg.Postgres.Port, cfg.Postgres.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}

	// Repositories
	blogRepo := repository.NewBlogRepository(db)
	authorRepo := repository.NewAuthorRepository(db)

	// Migrate
	if err := blogRepo.Migrate(); err != nil {
		log.Fatalf("failed to migrate Blog table: %v", err)
	}
	if err := authorRepo.Migrate(); err != nil {
		log.Fatalf("failed to migrate Author table: %v", err)
	}

	// RabbitMQ
	mqClient, err := rabbitmq.New(
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Exchange,
		cfg.RabbitMQ.ExchangeType,
	)
	if err != nil {
		log.Fatalf("failed to connect RabbitMQ: %v", err)
	}

	// Usecase
	blogUsecase := usecase.NewBlogUsecase(blogRepo, authorRepo, mqClient)

	// gRPC server
	grpcServer := grpc.NewServer()
	blogGRPCHandler := grpcHandler.NewBlogHandler(blogUsecase)
	blogpb.RegisterBlogServiceServer(grpcServer, blogGRPCHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", cfg.BlogService.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen gRPC: %v", err)
	}
	go func() {
		log.Printf("Blog gRPC listening at %s", cfg.BlogService.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC failed: %v", err)
		}
	}()

	// HTTP server
	mux := http.NewServeMux()
	blogHTTPHandler := httpHandler.NewBolgHandler(blogUsecase)

	// JWT middleware only for extracting user_id from request context
	jwtMiddleware := httpHandler.NewAuthMiddleware(nil) // Blog service doesn't generate/validate tokens itself

	mux.Handle("/blog/create", jwtMiddleware.RequireAuth(http.HandlerFunc(blogHTTPHandler.CreatePost)))

	httpServer := &http.Server{
		Addr:    cfg.BlogService.HTTPPort,
		Handler: mux,
	}

	go func() {
		log.Printf("Blog HTTP listening at %s", cfg.BlogService.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP failed: %v", err)
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	log.Println("Shutting down Blog service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown failed: %v", err)
	}

	grpcServer.GracefulStop()
	log.Println("Blog Service stopped")
}
