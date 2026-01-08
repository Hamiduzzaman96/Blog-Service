package grpc

import (
	"context"

	"github.com/Hamiduzzaman96/Blog-Service/internal/usecase"
	"github.com/Hamiduzzaman96/Blog-Service/proto/notificationpb"
)

type NotificationHandler struct {
	notificationpb.UnimplementedNotificationServiceServer
	usecase *usecase.NotificationUsecase
}

func NewNotificationHandler(u *usecase.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{usecase: u}
}

func (h *NotificationHandler) SendNotification(ctx context.Context, req *notificationpb.NotificationRequest) (*notificationpb.NotificationResponse, error) {
	err := h.usecase.Send(uint(req.UserId), req.Message)
	if err != nil {
		return nil, err
	}

	return &notificationpb.NotificationResponse{Success: true}, nil
}
