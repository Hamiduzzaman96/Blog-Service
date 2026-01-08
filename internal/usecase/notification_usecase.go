package usecase

import (
	"github.com/Hamiduzzaman96/Blog-Service/internal/domain"
	"github.com/Hamiduzzaman96/Blog-Service/internal/repository"
)

type NotificationUsecase struct {
	notificationRepo *repository.NotificationRepository
}

func NewNotificationUsecase(notificatioRepo *repository.NotificationRepository) *NotificationUsecase {
	return &NotificationUsecase{notificationRepo: notificatioRepo}
}

func (n *NotificationUsecase) Send(userID uint, message string) error {
	notification := &domain.Notification{
		UserID:  userID,
		Message: message,
		Sent:    false,
	}

	_, err := n.notificationRepo.Create(notification)
	return err
}
