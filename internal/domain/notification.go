package domain

type Notification struct {
	ID      uint
	UserID  uint
	Message string
	Sent    bool
}

func NewNotification(userID uint, message string) *Notification {
	return &Notification{
		UserID:  userID,
		Message: message,
		Sent:    false,
	}
}
