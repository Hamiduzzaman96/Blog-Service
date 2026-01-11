package repository

import (
	"errors"

	"github.com/Hamiduzzaman96/Blog-Service/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Migrate() error {
	return r.db.AutoMigrate(&UserModel{})
}

// Mapper //DB ---> Domain
func toDomainUser(m *UserModel) *domain.User {
	return &domain.User{
		ID:    m.ID,
		Email: m.Email,
		Role:  m.Role,
	}
}

// Mapper // Domain ---> DB
func toModelUser(u *domain.User, password string) *UserModel {
	return &UserModel{
		ID:       u.ID,
		Email:    u.Email,
		Role:     u.Role,
		Password: password,
	}
}

// CRUD
func (r *UserRepository) Create(u *domain.User, password string) (*domain.User, error) {
	m := toModelUser(u, password)
	if err := r.db.Create(&m).Error; err != nil {
		return nil, err
	}
	u.ID = m.ID
	return u, nil
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var m UserModel // local user model declare
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return toDomainUser(&m), nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var m UserModel
	if err := r.db.Where("email=?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return toDomainUser(&m), nil
}

func (r *UserRepository) Update(u *domain.User) error {
	return r.db.Model(&UserModel{}).Where("id = ?", u.ID).Updates(map[string]any{
		"email": u.Email,
		"role":  u.Role,
	}).Error
}
