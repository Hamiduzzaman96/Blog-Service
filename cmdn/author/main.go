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
	"github.com/Hamiduzzaman96/Blog-Service/pkg/jwt"
	"github.com/Hamiduzzaman96/Blog-Service/proto/authorpb"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load Config
	cfg := config.Load()

	// Signal context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// PostgreSQL Connection
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DB,
		cfg.Postgres.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	authorRepo := repository.NewAuthorRepository(db)

	// Auto migrate
	if err := userRepo.Migrate(); err != nil {
		log.Fatalf("failed to migrate User table: %v", err)
	}
	if err := authorRepo.Migrate(); err != nil {
		log.Fatalf("failed to migrate Author table: %v", err)
	}

	// JWT Service
	jwtSvc := jwt.NewService(cfg.JWT.Secret, cfg.JWT.AccessTokenExp, cfg.JWT.RefreshTokenExp)

	// Usecases
	authorUsecase := usecase.NewAuthorUsecase(userRepo, authorRepo)

	// gRPC Server
	grpcServer := grpc.NewServer()
	authorGRPCHandler := grpcHandler.NewAuthorHandler(authorUsecase)
	authorpb.RegisterAuthorServiceServer(grpcServer, authorGRPCHandler)

	lis, err := net.Listen("tcp", cfg.AuthorService.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen gRPC: %v", err)
	}
	go func() {
		log.Printf("Author gRPC listening at %s", cfg.AuthorService.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC failed: %v", err)
		}
	}()

	// HTTP Server
	mux := http.NewServeMux()
	authorHTTPHandler := httpHandler.NewAuthorHandler(authorUsecase)
	authMiddleware := httpHandler.NewAuthMiddleware(jwtSvc)

	// Become Author route with auth middleware
	mux.Handle("/become-author", authMiddleware.RequireAuth(http.HandlerFunc(authorHTTPHandler.BecomeAuthor)))

	httpServer := &http.Server{
		Addr:    cfg.AuthorService.HTTPPort,
		Handler: mux,
	}

	go func() {
		log.Printf("Author HTTP listening at %s", cfg.AuthorService.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP failed: %v", err)
		}
	}()

	// Wait for termination
	<-ctx.Done()
	stop()
	log.Println("Shutting down Author service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown failed: %v", err)
	}

	grpcServer.GracefulStop()
	log.Println("Author service stopped")
}
