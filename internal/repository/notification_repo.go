package repository

import (
	"github.com/Hamiduzzaman96/Blog-Service/internal/domain"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Migrate() error {
	return r.db.AutoMigrate(&NotificationModel{})
}

// MAPPPERS

func notificationModelToDomain(m *NotificationModel) *domain.Notification {
	return &domain.Notification{
		ID:      m.ID,
		UserID:  m.UserID,
		Message: m.Message,
		Sent:    m.Sent,
	}
}

func notificationDomainToModel(n *domain.Notification) *NotificationModel {
	return &NotificationModel{
		ID:      n.ID,
		UserID:  n.UserID,
		Message: n.Message,
		Sent:    n.Sent,
	}
}

// CRUD
func (r *NotificationRepository) Create(n *domain.Notification) (*domain.Notification, error) {
	m := notificationDomainToModel(n)

	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	n.ID = m.ID
	return n, nil
}

func (r *NotificationRepository) MarkAsSent(id uint) error {
	return r.db.Model(&NotificationModel{}).
		Where("id = ?", id).
		Update("sent", true).Error
}
