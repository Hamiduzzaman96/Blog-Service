package repository

import "time"

type UserModel struct {
	ID       uint   `gorm:"primarykey;autoIncrement"`
	Email    string `gorm:"unique; not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"` // User // Author

}

type AuthorModel struct {
	ID     uint `gorm:"primarykey;autoIncrement"`
	UserID uint `gorm:"unique; not null"`
}

type BlogModel struct {
	ID        uint   `gorm:"primarykey;autoIncrement"`
	AuthorID  uint   `gorm:"not null"`
	Title     string `gorm:"not null"`
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NotificationModel struct {
	ID      uint   `gorm:"primarykey; autoIncrement"`
	UserID  uint   `gorm:"not null"`
	Message string `gorm:"type:text"`
	Sent    bool
}
