package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hamiduzzaman96/Blog-Service/config"
	pkgJWT "github.com/Hamiduzzaman96/Blog-Service/pkg/jwt"
	pkgRedis "github.com/Hamiduzzaman96/Blog-Service/pkg/redis"
	"github.com/Hamiduzzaman96/Blog-Service/proto/userpb"

	grpcHandler "github.com/Hamiduzzaman96/Blog-Service/internal/handler/grpc"
	httpHandler "github.com/Hamiduzzaman96/Blog-Service/internal/handler/http"
	"github.com/Hamiduzzaman96/Blog-Service/internal/repository"
	"github.com/Hamiduzzaman96/Blog-Service/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dsn := "host=" + cfg.Postgres.Host +
		" user=" + cfg.Postgres.User +
		" password=" + cfg.Postgres.Password +
		" dbname=" + cfg.Postgres.DB +
		" port=" + cfg.Postgres.Port +
		" sslmode=" + cfg.Postgres.SSLMode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	if err := userRepo.Migrate(); err != nil {
		log.Fatalf("failed to migrate User table: %v", err)
	}

	redisClient, err := pkgRedis.New(cfg.Redis.Host+":"+cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB, 2*time.Hour)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
	defer redisClient.Close()

	jwtSvc := pkgJWT.NewService(cfg.JWT.Secret, cfg.JWT.AccessTokenExp, cfg.JWT.RefreshTokenExp)

	userUsecase := usecase.NewUserUsecase(userRepo, jwtSvc, redisClient)

	grpcServer := grpc.NewServer()
	userGRPCHandler := grpcHandler.NewUserHandler(userUsecase)
	userpb.RegisterUserServiceServer(grpcServer, userGRPCHandler)
	reflection.Register(grpcServer)

	grpcLis, err := net.Listen("tcp", cfg.UserService.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen gRPC: %v", err)
	}
	go func() {
		log.Printf("User gRPC listening at %s", cfg.UserService.GRPCPort)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("gRPC failed: %v", err)
		}
	}()

	userHTTPHandler := httpHandler.NewUserHandler(userUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/register", userHTTPHandler.Register)
	mux.HandleFunc("/login", userHTTPHandler.Login)

	httpServer := &http.Server{
		Addr:    cfg.UserService.HTTPPort,
		Handler: mux,
	}
	go func() {
		log.Printf("User HTTP listening at %s", cfg.UserService.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP failed: %v", err)
		}
	}()

	<-ctx.Done()
	stop()
	log.Println("Shutting down User service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown failed: %v", err)
	}
	grpcServer.GracefulStop()
	log.Println("User service stopped")
}
