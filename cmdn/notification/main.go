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
	"github.com/Hamiduzzaman96/Blog-Service/proto/notificationpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password,
		cfg.Postgres.DB, cfg.Postgres.Port, cfg.Postgres.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}

	notifRepo := repository.NewNotificationRepository(db)

	if err := notifRepo.Migrate(); err != nil {
		log.Fatalf("failed to migrate Notification table: %v", err)
	}

	notifUsecase := usecase.NewNotificationUsecase(notifRepo)

	grpcServer := grpc.NewServer()
	notifGRPCHandler := grpcHandler.NewNotificationHandler(notifUsecase)
	notificationpb.RegisterNotificationServiceServer(grpcServer, notifGRPCHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", cfg.NotificationService.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen gRPC port: %v", err)
	}
	go func() {
		log.Printf("Notification gRPC listening at %s", cfg.NotificationService.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC failed: %v", err)
		}
	}()

	mux := http.NewServeMux()
	notifHTTPHandler := httpHandler.NewNotificationHandler(notifUsecase)

	mux.Handle("/send-notification", http.HandlerFunc(notifHTTPHandler.SendNotification))

	httpServer := &http.Server{
		Addr:    cfg.NotificationService.HTTPPort,
		Handler: mux,
	}

	go func() {
		log.Printf("Notification HTTP listening at %s", cfg.NotificationService.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP failed: %v", err)
		}
	}()

	<-ctx.Done()
	cancel()
	log.Println("Shutting down Notification Service...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown failed: %v", err)
	}

	grpcServer.GracefulStop()
	log.Println("Notification Service stopped")
}
